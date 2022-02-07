package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()                                       // 创建实例
	app.HandleDir("/public", "./fronted/web/public")        //设置模版目标
	app.HandleDir("/html", "./fronted/web/htmlProductShow") // 访问生成好的静态文件
	app.Run(
		iris.Addr("localhost:8083"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
