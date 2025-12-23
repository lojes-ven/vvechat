package service

import (
	"vvechat/internal/model"
	"vvechat/pkg/infra"

	"gorm.io/gorm"
)

func FriendshipList(userID uint64) ([]model.FriendshipListResp, error) {
	var resp []model.FriendshipListResp

	res := infra.GetDB().
		Model(&model.Friendship{}).
		Where("user_id = ?", userID).
		Find(&resp)

	if res.Error != nil {
		return nil, res.Error
	}

	return resp, nil
}

// createFriendship 给两个人（id主键）创建出好友关系
func createFriendship(tx *gorm.DB, userID, friendID uint64) error {
	var friendName, userName string

	err := tx.Table("users").
		Select("name").
		Where("id = ?", friendID).
		Row().
		Scan(&friendName)

	if err != nil {
		return err
	}

	err = tx.Table("users").
		Select("name").
		Where("id = ?", userID).
		Row().
		Scan(&userName)

	if err != nil {
		return err
	}

	res := tx.Model(&model.Friendship{}).
		Create(model.NewFriendship(userID, friendID, friendName))
	if res.Error != nil {
		return res.Error
	}

	res = tx.Model(&model.Friendship{}).
		Create(model.NewFriendship(friendID, userID, userName))
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteFriendship(userID, friendID uint64) error {
	res := infra.GetDB().
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Delete(&model.Friendship{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
