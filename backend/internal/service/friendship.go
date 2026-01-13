package service

import (
	"errors"
	"log"
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
		log.Println(res.Error)
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
		log.Println(err)
		return err
	}

	err = tx.Table("users").
		Select("name").
		Where("id = ?", userID).
		Row().
		Scan(&userName)

	if err != nil {
		log.Println(err)
		return err
	}

	res := tx.Model(&model.Friendship{}).
		Create(model.NewFriendship(userID, friendID, friendName))
	if res.Error != nil {
		log.Println(res.Error)
		return res.Error
	}

	res = tx.Model(&model.Friendship{}).
		Create(model.NewFriendship(friendID, userID, userName))
	if res.Error != nil {
		log.Println(res.Error)
		return res.Error
	}
	return nil
}

func DeleteFriendship(userID, friendID uint64) error {
	db := infra.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		res := db.Where("user_id = ? AND friend_id = ?", userID, friendID).
			Delete(&model.Friendship{})
		if res.Error != nil {
			log.Println(res.Error)
			return res.Error
		}
		if res.RowsAffected == 0 {
			log.Println("删除好友操作影响了0行表")
			return gorm.ErrRecordNotFound
		}

		res = db.Where("user_id = ? AND friend_id = ?", friendID, userID).
			Delete(&model.Friendship{})
		if res.Error != nil {
			log.Println(res.Error)
			return res.Error
		}
		if res.RowsAffected == 0 {
			log.Println("删除好友操作影响了0行表")
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func ReviseRemark(userID, friendID uint64, remark string) error {
	db := infra.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Friendship{}).
			Where("user_id = ? AND friend_id = ?", userID, friendID).
			Update("friend_remark", remark)
		if res.Error != nil {
			log.Println(res.Error)
			return errors.New("服务器错误")
		}
		if res.RowsAffected == 0 {
			log.Println("修改备注操作影响了0行表")
			return errors.New("服务器错误")
		}

		res = tx.Model(&model.ConversationUser{}).
			Where("user_id = ? AND conversation_id = ?", userID, friendID).
			Update("remark", remark)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			log.Println(res.Error)
			return errors.New("服务器错误")
		}
		return nil
	})
}
