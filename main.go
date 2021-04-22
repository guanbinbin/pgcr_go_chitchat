package main

import (
	"github.com/pygosuperman/pgcr_go_chitchat/route"
	"net/http"
	"time"
)

func main() {
	route.P("ChitChat", route.Version(), "started at", route.Config.Address)

	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(route.Config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// all route patterns matched here
	// route handler functions defined in other files
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
		Addr:           route.Config.Address,
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
