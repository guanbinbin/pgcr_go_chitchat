package main

import (
	"net/http"
	"time"

	"github.com/pygosuperman/pgcr_go_chitchat/route"
)

func main() {
	// 打印一句话，表示程序开始运行
	route.P("pgcr_go_chitchat", route.Version(), "开始运行，地址是：", route.Config.Address)

	// 这是一个多路复用器，作为所有请求的入口
	// 每个请求都会被传递到这里
	mux := http.NewServeMux()
	// 静态资源目录，通过配置文件可以进行配置
	files := http.FileServer(http.Dir(route.Config.Static))
	// 静态资源的路径
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// 所有的路由都在这里匹配
	// 具体的路由函数定义在其他的文件中
	//

	// index
	mux.HandleFunc("/", route.Index)
	// error
	mux.HandleFunc("/err", route.Err)

	// defined in route_auth.go
	mux.HandleFunc("/login", route.Login)
	mux.HandleFunc("/logout", route.Logout)
	mux.HandleFunc("/signup", route.Signup)
	mux.HandleFunc("/signup_account", route.SignupAccount)
	mux.HandleFunc("/authenticate", route.Authenticate)

	// defined in route_thread.go
	mux.HandleFunc("/thread/new", route.NewThread)
	mux.HandleFunc("/thread/create", route.CreateThread)
	mux.HandleFunc("/thread/post", route.PostThread)
	mux.HandleFunc("/thread/read", route.ReadThread)

	// starting up the server
	server := &http.Server{
		Addr: route.Config.Address,
		// 指定处理器为一个多路复用器对象
		Handler:        mux,
		ReadTimeout:    time.Duration(route.Config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(route.Config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	// server := &http.Server{
	// 	Addr:           "0.0.0.0:8080",
	// 	Handler:        mux,
	// 	ReadTimeout:    time.Duration(10 * int64(time.Second)),
	// 	WriteTimeout:   time.Duration(600 * int64(time.Second)),
	// 	MaxHeaderBytes: 1 << 20,
	// }
	server.ListenAndServe()
}
