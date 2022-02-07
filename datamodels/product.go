package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" form:"ID" column:"ID"`
	ProductName  string `json:"ProductName" sql:"productName" form:"ProductName" gorm:"column:productName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" form:"ProductNum" gorm:"column:productNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" form:"ProductImage" gorm:"column:productImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" form:"ProductUrl" gorm:"column:productUrl"`
}
