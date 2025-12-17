package model

type RegisterReq struct {
	Name        string `json:"name" binding:"required,min=1,max=64"`
	Password    string `json:"password" binding:"required,min=6,max=72"`
	PhoneNumber string `json:"phone_number" binding:"required,len=11,numeric"`
}

type LoginByUidReq struct {
	Password string `json:"password" binding:"required,min=6,max=72"`
	Uid      string `json:"uid" binding:"required,min=1,max=20"`
}

type LoginByPhoneReq struct {
	PhoneNumber string `json:"phone_number" binding:"required,len=11,numeric"`
	Password    string `json:"password" binding:"required,min=6,max=72"`
}
