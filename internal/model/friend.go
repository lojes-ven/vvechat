package model

import (
	"gorm.io/gorm"
)

// Friendship 好友关系表
type Friendship struct {
	gorm.Model
	UserID   uint64 `gorm:"type:bigint;not null;index:idx_user_friend,unique"`
	FriendID uint64 `gorm:"type:bigint;not null;index:idx_user_friend,unique"`
}

// FriendshipRequest 好友申请列表
type FriendshipRequest struct {
	gorm.Model
	SenderID            uint64 `gorm:"type:bigint;not null;index:idx_sender_receiver,unique"`
	SenderName          string `gorm:"type:varchar(64);not null;"`
	ReceiverID          uint64 `gorm:"type:bigint;not null;index:idx_sender_receiver,unique"`
	VerificationMessage string `gorm:"type:varchar(128)"`
	Status              string `gorm:"type:varchar(16);not null;check:status IN ('pending','accepted','rejected','canceled')"`
}

func NewFriendship(id1, id2 uint64) *Friendship {
	return &Friendship{
		UserID:   id1,
		FriendID: id2,
	}
}

func NewFriendshipRequest(senderID, receiverID uint64, msg string, senderName string) *FriendshipRequest {
	return &FriendshipRequest{
		SenderID:            senderID,
		SenderName:          senderName,
		ReceiverID:          receiverID,
		VerificationMessage: msg,
		Status:              "pending",
	}
}
