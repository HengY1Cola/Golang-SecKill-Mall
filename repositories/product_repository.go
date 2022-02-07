package repositories

import (
	"Iris/common"
	"Iris/datamodels"
	"gorm.io/gorm"
)

//  todo 先开发对应的接口

type IProduct interface {
	Conn() error //连接数据
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
	SubProductNum(productID int64) error
}

//todo 实现定义的接口

type ProductManager struct {
	table     string
	mysqlConn *gorm.DB
}

// NewProductManager 创建构造函数
func NewProductManager(table string, db *gorm.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

// todo 实现接口功能

// Conn 数据连接
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewGormMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	return
}

// Insert 插入
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	// todo 判断连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	// todo 进行构造
	productStruct := &datamodels.Product{
		ProductName:  product.ProductName,
		ProductNum:   product.ProductNum,
		ProductImage: product.ProductImage,
		ProductUrl:   product.ProductUrl,
	}
	// todo 开始插入
	result := p.mysqlConn.Table("product").Create(&productStruct)
	if result.Error != nil {
		return 0, result.Error
	}
	return productStruct.ID, nil
}

// Delete 商品的删除
func (p *ProductManager) Delete(productID int64) bool {
	// todo 判断连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}
	// todo 构造并删除
	product := &datamodels.Product{}
	res := p.mysqlConn.Table("product").Delete(&product, productID)
	// todo 判断
	if res.Error != nil {
		return false
	}
	return true
}

// Update 商品的更新
func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}

	productStruct := &datamodels.Product{}
	res := p.mysqlConn.Table("product").First(&productStruct, product.ID)
	if res.Error != nil {
		return res.Error
	}

	productStruct.ProductName = product.ProductName
	productStruct.ProductNum = product.ProductNum
	productStruct.ProductImage = product.ProductImage
	productStruct.ProductUrl = product.ProductUrl
	res = p.mysqlConn.Table("product").Save(&productStruct)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// SelectByKey 根据商品ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	if err := p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}

	productResult = &datamodels.Product{}
	res := p.mysqlConn.Table("product").First(&productResult, productID)
	if res.Error != nil {
		return &datamodels.Product{}, res.Error
	}
	return
}

// SelectAll 获取所有商品
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error) {
	if err := p.Conn(); err != nil {
		return []*datamodels.Product{}, err
	}

	var product []*datamodels.Product
	p.mysqlConn.Table("product").Find(&product) // 结构体的复数寻找表名
	for _, v := range product {
		productArray = append(productArray, v)
	}
	return
}

func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}

	productStruct := &datamodels.Product{}
	res := p.mysqlConn.Table("product").First(&productStruct, productID)
	if res.Error != nil {
		return res.Error
	}
	productStruct.ProductNum = productStruct.ProductNum - 1
	res = p.mysqlConn.Table("product").Save(&productStruct)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
