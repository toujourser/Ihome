package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"net/http"
	"sss/IhomeWeb/handler"
	_ "sss/IhomeWeb/models"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":8008"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("===-------------------------")

	rou := httprouter.New()
	rou.NotFound = http.FileServer(http.Dir("html"))
	// 获取地区信息
	rou.GET("/api/v1.0/areas", handler.GetArea)
	// 获取首页轮播图
	rou.GET("/api/v1.0/house/index", handler.GetIndex)
	// 获取session
	rou.GET("/api/v1.0/session", handler.GetSession)
	// 获取验证码图片
	rou.GET("/api/v1.0/imagecode/:uuid", handler.GetImages)
	// 获取短信验证码
	rou.GET("/api/v1.0/smscode/:mobile", handler.GetSmsCode)
	//
	rou.POST("/api/v1.0/users", handler.PostRet)
	service.Handle("/", rou)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
