package route

import (
	"net/http"

	"github.com/pygosuperman/pgcr_go_chitchat/data"
)

// GET /err?msg=
// 展示错误页面
func Err(writer http.ResponseWriter, request *http.Request) {
	// 获取查询参数
	vals := request.URL.Query()
	// 校验是否登录
	_, err := session(writer, request)
	if err != nil {
		// 没有登录，展示公共的错误页面
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		// 登录了，展示个人的错误页面
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

// 首页的处理器函数
func Index(writer http.ResponseWriter, request *http.Request) {
	// 获取所有的帖子
	threads, err := data.Threads()
	if err != nil {
		// 打印错误消息
		error_message(writer, request, "获取帖子失败")
	} else {
		// 检查用户是否登录
		_, err := session(writer, request)
		if err != nil {
			// 生成公共模板
			generateHTML(writer, threads, "layout", "public.navbar", "index")
		} else {
			// 生成私有模板
			generateHTML(writer, threads, "layout", "private.navbar", "index")
		}
	}
}
