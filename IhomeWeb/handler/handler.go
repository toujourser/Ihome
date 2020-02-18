package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"image"
	"image/png"
	"net/http"
	"regexp"
	go_micro_srv_GetArea "sss/GetArea/proto/GetArea"
	go_micro_srv_GetImageCd "sss/GetImageCd/proto/GetImageCd"
	go_micro_srv_GetSmscd "sss/GetSmscd/proto/GetSmscd"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	go_micro_srv_PostRet "sss/PostRet/proto/PostRet"
)

func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取地区请求客户端 url:api/v1.0/areas")
	server := grpc.NewService(micro.Name("go.micro.web.IhomeWeb"))
	server.Init()

	// 调用服务返回句柄
	client := go_micro_srv_GetArea.NewGetAreaService("go.micro.srv.GetArea", server.Client())

	resp, err := client.Call(context.TODO(), &go_micro_srv_GetArea.Request{})
	if err != nil {
		beego.Info("err: == ", err)
		beego.Info("resp: == ", resp)
		http.Error(w, err.Error(), 500)
		return
	}

	area_list := []models.Area{}

	for _, value := range resp.Data {
		tmp := models.Area{
			Id:   int(value.Aid),
			Name: value.Aname,
		}
		area_list = append(area_list, tmp)
	}

	response := map[string]interface{}{
		"errno":  resp.Error,
		"errmsg": resp.Errmsg,
		"data":   area_list,
	}

	// 回传数据的时候需要设置数据格式
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func GetIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("获取首页轮播图 url:/api/v1.0/house/index")
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetSession(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("获取Session url:/api/v1.0/session")
	response := map[string]interface{}{
		"errno":  utils.RECODE_SESSIONERR,
		"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("获取首页轮播图 url: /api/v1.0/imagecode/:uuid")
	uuid := p.ByName("uuid")
	fmt.Println("uuid: == ", uuid)

	server := grpc.NewService()
	server.Init()

	client := go_micro_srv_GetImageCd.NewGetImageCdService("go.micro.srv.GetImageCd", server.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_GetImageCd.Request{
		Uuid: uuid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var img image.RGBA
	for _, value := range resp.Pix {
		img.Pix = append(img.Pix, uint8(value))
	}
	img.Stride = int(resp.Stride)
	img.Rect.Min.X = int(resp.Min.X)
	img.Rect.Min.Y = int(resp.Min.Y)
	img.Rect.Max.X = int(resp.Max.X)
	img.Rect.Max.Y = int(resp.Max.Y)

	var image captcha.Image
	image.RGBA = &img
	fmt.Println(image)

	png.Encode(w, image)

}

func GetSmsCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("获取短信验证码 url: /api/v1.0/smscode/:mobile")
	mobile := p.ByName("mobile")
	// 正则表达式 手机号
	mobileReg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	bl := mobileReg.MatchString(mobile)
	if !bl {
		beego.Info("+++++", bl)
		utils.Response(w, utils.RECODE_MOBILEERR, utils.RecodeText(utils.RECODE_MOBILEERR), nil)
		return

	}
	beego.Info("====================")

	imageStr := r.URL.Query()["text"][0]
	uuid := r.URL.Query()["id"][0]

	server := grpc.NewService()
	server.Init()

	client := go_micro_srv_GetSmscd.NewGetSmscdService("go.micro.srv.GetSmscd", server.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_GetSmscd.Request{
		Mobile:   mobile,
		Uuid:     uuid,
		ImageStr: imageStr,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := utils.Response(w, resp.Error, resp.Errmsg, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func PostRet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("用户注册 url: api/v1.0/users")
	var reqParams = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&reqParams); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	beego.Info(reqParams, "===========")
	if reqParams["mobile"] == "" || reqParams["password"] == "" || reqParams["sms_code"] == "" {
		utils.Response(w, utils.RECODE_PARAMERR, utils.RecodeText(utils.RECODE_PARAMERR), nil)
		return
	}

	server := grpc.NewService()
	server.Init()

	client := go_micro_srv_PostRet.NewPostRetService("go.micro.srv.PostRet", server.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_PostRet.Request{
		Mobile:   reqParams["mobile"],
		Password: reqParams["password"],
		SmsCode:  reqParams["sms_code"],
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// 读取cookie 统一 "userLogin"
	cookie, err := r.Cookie("userLogin")
	if err != nil || "" == cookie.Value {
		cookie := &http.Cookie{
			Name:   "userLogin",
			Value:  resp.SessionId,
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(w, cookie)
	}
	if err := utils.Response(w, resp.Erron, resp.Errmsg, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
