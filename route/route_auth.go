package route

import (
	"net/http"

	"github.com/pygosuperman/pgcr_go_chitchat/data"
)

// GET /login
// Show the Login page
func Login(writer http.ResponseWriter, request *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(writer, nil)
}

// GET /signup
// Show the Signup page
func Signup(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// Create the user account
func SignupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	user := data.User{
		Name:     request.PostFormValue("name"),
		Email:    request.PostFormValue("email"),
		Password: request.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	http.Redirect(writer, request, "/login", http.StatusFound)
}

// POST /authenticate
// Authenticate the user given the email and password
// 校验用户输入的邮箱和密码
func Authenticate(writer http.ResponseWriter, request *http.Request) {
	// 解析表单数据
	err := request.ParseForm()
	if err != nil {
		danger(err, "未能正确解析表单数据")
	}
	// 根据邮箱查找用户
	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		danger(err, "未能找到用户")
	}
	// 校验用户输入的密码加密以后，是否与数据库中查询出的用户的密码匹配
	if user.Password == data.Encrypt(request.PostFormValue("password")) {
		// 创建一个session对象
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		// 创建一个cookie对象
		cookie := http.Cookie{
			Name: "_cookie",
			// 从session中取出uuid，作为cookie的值
			Value: session.Uuid,
			// 不允许客户端使用JavaScript取出操作
			HttpOnly: true,
		}
		// 设置cookie
		http.SetCookie(writer, &cookie)
		// 重定向到首页
		http.Redirect(writer, request, "/", http.StatusFound)
	} else {
		// 重定向到登录页面
		http.Redirect(writer, request, "/login", http.StatusFound)
	}

}

// GET /logout
// Logs the user out
func Logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		warning(err, "Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}
