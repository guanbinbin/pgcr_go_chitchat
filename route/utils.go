package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pygosuperman/pgcr_go_chitchat/data"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var Config Configuration
var logger *log.Logger

// Convenience function for printing to stdout
// 标准输出的便捷函数
func P(a ...interface{}) {
	fmt.Println(a)
}

func init() {
	loadConfig()
	file, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

// Convenience function to redirect to the error message page
func error_message(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

// 检查用户是否已经登录
func session(writer http.ResponseWriter, request *http.Request) (sess data.Session, err error) {
	// 获取cookie数据
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		// 创建一个session对象
		sess = data.Session{Uuid: cookie.Value}
		// 检查session是否存在
		if ok, _ := sess.Check(); !ok {
			// 不存在，报错
			err = errors.New("session不存在，用户未登录") //文字内容不应该以大写字母开头或者标点符号结尾。
		}
	}
	return
}

// 解析HTML模板
// 通过文件名称的列表，得到一个模板对象
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	// template.Must：用于捕捉语法分析过程中可能会产生的错误
	t = template.Must(t.ParseFiles(files...))
	return
}

// 生成HTML的函数
func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	// 追加文件列表
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	// 编译模板
	templates := template.Must(template.ParseFiles(files...))
	// 执行模板，传递数据
	templates.ExecuteTemplate(writer, "layout", data)
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}

// Version
func Version() string {
	return "0.1"
}
