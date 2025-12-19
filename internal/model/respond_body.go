package model

import "time"

// 用户信息返回体
type UserInfoResp struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

// 好友信息返回体
type FriendInfoResp struct {
	ID     uint64 `json:"id"`
	Remark string `json:"remark"`
	Name   string `json:"name"`
	Uid    string `json:"uid"`
}

// 陌生人信息返回体
type StrangerInfoResp struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// 刷新token操作返回体
type TokenResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint64 `json:"expires_in"`
}

// 登陆操作返回体
type LoginResp struct {
	UserInfo   UserInfoResp `json:"user_info"`
	TokenClass TokenResp    `json:"token_class"`
}

// 好友申请列表返回体
type FriendRequestListResp struct {
	RequestID           uint64    `gorm:"column:id" json:"request_id"`
	SenderID            uint64    `json:"sender_id"`
	SenderName          string    `json:"sender_name"`
	VerificationMessage string    `json:"verification_message"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}
