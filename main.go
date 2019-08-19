package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	mvc.Configure(app.Party("/basic"), basicMVC)
	app.Run(iris.Addr(":8080"))
}

func basicMVC(app *mvc.Application) {

	//使用中间件
	app.Router.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Path:%s", ctx.Path())
		ctx.Next()
	})

	//注册绑定到控制器上的依赖，依赖可以是函数：参数为iris.Context、返回值是单纯值，也可以是一个结构体值
	//Register dependencies which will be binding to the controller(s),
	//can be either a function which accepts an iris.Context and returns a single value (dynamic binding)
	//or a static struct value (service).
	app.Register(
		sessions.New(sessions.Config{}).Start,
		&prefixedLogger{prefix: "DEV"},
	)

	//子控制器(basicController、basicSubController)有父控制器(*mvc.Application)中依赖的clone，故所有控制器访问同一个会话
	app.Handle(new(basicController))
	app.Party("/sub").Handle(new(basicSubController))
}

type LoggerService interface {
	Log(string)
}
type prefixedLogger struct {
	prefix string
}

func (s *prefixedLogger) Log(msg string) {
	fmt.Printf("%s: %s\n", s.prefix, msg)
}

type basicController struct {
	Logger  LoggerService
	Session *sessions.Session
}

func (c *basicController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/custom", "Custom")
	//b.HandleMany("GET", "/custom custom2", "Custom") HandleMany 2019-07-11后新增的方法
}

func (c basicController) AfterActivation(a mvc.AfterActivation) {
	if a.Singleton() {
		panic("basicController should be stateless,a request-scoped,we have a 'Session' which depends on the context.")
	}
}

func (c *basicController) Get() string {
	count := c.Session.Increment("count", 1)
	body := fmt.Sprintf("Hello from basicController\nTotal vistors from you: %d", count)
	c.Logger.Log(body)
	return body
}

func (c *basicController) Custom() string {
	return "custom"
}

type basicSubController struct {
	Session *sessions.Session
}

func (c basicSubController) Get() string {
	count := c.Session.GetIntDefault("count", 1)
	return fmt.Sprintf("Hello from basicSubController.\nRead-only visits count: %d", count)
}
