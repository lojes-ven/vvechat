package service

import (
	"log"
	"vvechat/internal/model"
	"vvechat/pkg/infra"
)

// RefreshToken 刷新Token
func RefreshToken(id uint64) (*model.TokenResp, error) {
	return NewTokenResp(id)
}

// ReviseUid 修改微信号
func ReviseUid(id uint64, newUid string) error {
	return infra.GetDB().
		Model(&model.User{}).
		Where("id = ?", id).
		Update("uid", newUid).
		Error
}

// FriendInfo 查看好友信息
func FriendInfo(userID, friendID uint64) (*model.FriendInfoResp, error) {
	var resp model.FriendInfoResp
	resp.ID = friendID
	db := infra.GetDB()

	res := db.Raw(`SELECT f.friend_remark, u.uid, u.name
		FROM friendships f
		JOIN users u 
		ON u.id = f.friend_id
		WHERE f.user_id = ? AND f.friend_id = ?
	`, userID, friendID).Scan(&resp)

	if res.Error != nil {
		log.Println(res.Error)
		return nil, res.Error
	}

	return &resp, nil
}

// StrangerInfo 查看陌生人信息
func StrangerInfo(userID, strangerID uint64) (*model.StrangerInfoResp, error) {
	var resp model.StrangerInfoResp
	resp.ID = strangerID
	db := infra.GetDB()

	res := db.Table("users").
		Select("name").
		Where("id = ?", strangerID).
		First(&resp)
	if res.Error != nil {
		log.Println(res.Error)
		return nil, res.Error
	}

	return &resp, nil
}
