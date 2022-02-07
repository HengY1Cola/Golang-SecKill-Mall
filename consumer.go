package main

import (
	"Iris/common"
	"Iris/rabbitmq"
	"Iris/repositories"
	"Iris/services"
	"log"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	gormDb, err := common.NewGormMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	// todo 创建数据库操作实例
	product := repositories.NewProductManager("product", gormDb)
	productService := services.NewProductService(product)
	order := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(order)
	// todo 消费消息队列
	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("Product")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
