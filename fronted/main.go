package main

import (
	"Iris/common"
	"Iris/fronted/middleware"
	"Iris/fronted/web/controllers"
	"Iris/rabbitmq"
	"Iris/repositories"
	"Iris/services"
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
)

func main() {
	// todo 基础设置
	app := iris.New()              // 创建实例
	app.Logger().SetLevel("debug") //设置模式
	// todo views设置
	template := iris.HTML("./fronted/web/views", ".html").
		Layout("shared/layout.html").
		Reload(true) // 注册模版
	app.RegisterView(template)
	app.HandleDir("/public", "./fronted/web/public")        //设置模版目标
	app.HandleDir("/html", "./fronted/web/htmlProductShow") // 访问生成好的静态文件
	// todo 异常页面处理
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	// todo 链接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	gormDb, err := common.NewGormMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//sess := sessions.New(sessions.Config{
	//	Cookie:  "setCookie",
	//	Expires: 60 * time.Minute,
	//})
	// todo 注册消息队列
	rabbitMq := rabbitmq.NewRabbitMQSimple("Product")
	// todo 注册控制器
	// 用户路由器
	userRepository := repositories.NewUserRepository("user", db)
	userService := services.NewService(userRepository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	//user.Register(userService, ctx, sess.Start)
	user.Register(userService, ctx)
	user.Handle(new(controllers.UserController))
	// 商品与订单路由器
	product := repositories.NewProductManager("product", gormDb)
	productService := services.NewProductService(product)
	order := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct)
	pro.Register(productService, orderService, ctx, rabbitMq)
	pro.Handle(new(controllers.ProductController))
	// todo 启动服务
	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
