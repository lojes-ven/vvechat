package service

import (
	"vvechat/internal/model"
	"vvechat/pkg/infra"

	"gorm.io/gorm"
)

// 发送好友申请操作
func SendFriendRequest(senderID, receiverID uint64, msg string, senderName string) error {
	if senderID == receiverID {
		return gorm.ErrInvalidData
	}
	if err := isPKExist(receiverID); err != nil {
		return err
	}

	return infra.GetDB().
		Model(&model.FriendshipRequest{}).
		Create(model.NewFriendshipRequest(senderID, receiverID, msg, senderName)).
		Error
}

// 加载好友申请列表操作
func FriendRequestList(receiverID uint64) ([]model.FriendRequestListResp, error) {
	respSlice := make([]model.FriendRequestListResp, 0)

	res := infra.GetDB().
		Model(&model.FriendshipRequest{}).
		Where("receiver_id = ?", receiverID).
		Order("created_at DESC").
		Find(&respSlice)
	if res.Error != nil {
		return nil, res.Error
	}

	return respSlice, nil
}
