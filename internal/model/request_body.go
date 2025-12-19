package model

// 注册请求体
type RegisterReq struct {
	Name        string `json:"name" binding:"required,min=1,max=64"`
	Password    string `json:"password" binding:"required,min=6,max=72"`
	PhoneNumber string `json:"phone_number" binding:"required,len=11,numeric"`
}

// 微信号登陆请求体
type LoginByUidReq struct {
	Password string `json:"password" binding:"required,min=6,max=72"`
	Uid      string `json:"uid" binding:"required,min=1,max=20"`
}

// 手机号登陆请求体
type LoginByPhoneReq struct {
	PhoneNumber string `json:"phone_number" binding:"required,len=11,numeric"`
	Password    string `json:"password" binding:"required,min=6,max=72"`
}

// 修改微信号请求体
type ReviseUidReq struct {
	NewUid string `json:"new_uid" binding:"required,min=1,max=20"`
}

// 添加好友请求体
type AddFriendReq struct {
	ReceiverID          uint64 `json:"receiver_id" binding:"required,gt=0"`
	SenderName          string `json:"sender_name" binding:"required,max=64"`
	VerificationMessage string `json:"verification_message" binding:"omitempty,max=128"`
}
