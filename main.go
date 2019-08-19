package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
)

func main() {
	app := newApp()
	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	// http://localhost:8080/custom_path
	app.Run(iris.Addr(":8080"))

}

func newApp() *iris.Application {
	app := iris.New()
	app.Use(recover.New())                      //可以从任何与http相关的panics中恢复
	app.Use(logger.New())                       //请求记录到终端
	mvc.New(app).Handle(new(ExampleController)) //控制器根路由路径"/"
	return app
}

type ExampleController struct{}

//Get service
//Method GET
//Resource:http://localhost:8080
func (c *ExampleController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

//遵循命名规范
// GetPing 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/ping
func (c *ExampleController) GetPing() string {
	return "pong"
}

//遵循命名规范
// GetHello 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/hello
func (c *ExampleController) GetHello() interface{} {
	return map[string]string{"message": "Hello Iris!"}
}

func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
	anyMiddlewareHere := func(ctx iris.Context) {
		ctx.Application().Logger().Warnf("Inside / custom_path")
		ctx.Next()
	}
	//自定义路径 不遵循命名规范
	b.Handle("GET", "/custom_path", "CustomHandlerWithoutFollowingTheNamingGuide", anyMiddlewareHere)
}

// 不遵循命名规范 BeforeActivation 中 b.Handle方法注册自定义路径
// CustomHandlerWithoutFollowingTheNamingGuide 服务
// 请求方法:   GET
// 请求资源路径: http://localhost:8080/custom_path
func (c *ExampleController) CustomHandlerWithoutFollowingTheNamingGuide() string {
	return "hello from the custom handler without following the naming guide"
}
