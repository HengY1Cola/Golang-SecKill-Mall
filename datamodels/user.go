package datamodels

type User struct {
	ID           int64  `json:"id" form:"ID" sql:"ID"`
	NickName     string `json:"nickName" form:"nickName" sql:"nickName"`
	UserName     string `json:"userName" form:"userName" sql:"userName"`
	HashPassword string `json:"-" form:"passWord" sql:"passWord"`
}
