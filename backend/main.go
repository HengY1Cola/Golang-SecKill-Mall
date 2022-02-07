package main

import (
	"Iris/backend/web/controllers"
	"Iris/common"
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
	template := iris.HTML("./backend/web/views", ".html").
		Layout("shared/layout.html").
		Reload(true) // 注册模版
	app.RegisterView(template)
	app.HandleDir("/assets", "./backend/web/assets") //设置模版目标
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
	// todo 注册控制器
	// 商品路由器
	productRepository := repositories.NewProductManager("product", gormDb)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))
	// 订单路由器
	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))
	// todo 启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
