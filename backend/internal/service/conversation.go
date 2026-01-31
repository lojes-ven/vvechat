package service

import (
	"errors"
	"log"

	"github.com/lojes7/inquire/internal/model"
	"github.com/lojes7/inquire/pkg/infra"
	"github.com/lojes7/inquire/pkg/utils"
	"gorm.io/gorm"
)

// StartPrivateConversation 发起私聊
// 会调用getPrivateConversationID以获取会话ID
// 最后返回会话ID
func StartPrivateConversation(userID, friendID uint64) (uint64, error) {
	db := infra.GetDB()

	// 找到 A 和 B 共同的 conversation_id
	conversationID, err := getPrivateConversationID(userID, friendID)
	if err != nil {
		return 0, err
	}

	// 第二步：不管用户有没有删除该会话，都更新该用户的 conversation_users 记录
	// 使该会话在可能被删除的情况下重新出现
	res := db.Model(&model.ConversationUser{}).
		Where("user_id = ? AND conversation_id = ?", userID, conversationID).
		Update("deleted_at", nil)

	if res.Error != nil {
		log.Println(res.Error)
		return 0, errors.New("服务器错误")
	}

	return conversationID, nil
}

// ChatHistoryList 加载聊天记录
func ChatHistoryList(userID, conversationID uint64) ([]model.ChatHistoryResp, error) {
	db := infra.GetDB()

	resp := make([]model.ChatHistoryResp, 0)

	sql := `SELECT m.id AS message_id, 
       		m.sender_id, 
       		u.name AS sender_name,
			m.status, 
			m.updated_at,
			CASE m.status
			WHEN ? OR ? THEN t.text
			WHEN ? THEN json_build_object(
               'file_name', f.file_name,
               'file_url', f.file_url,
               'file_size', f.file_size,
               'file_type', f.file_type
           )
			ELSE ''
			END AS content
			FROM messages m
			LEFT JOIN users u ON u.id = m.sender_id
			LEFT JOIN message_users mu ON mu.message_id = m.id AND mu.user_id = ? 
			LEFT JOIN texts t ON t.message_id = m.id
			LEFT JOIN files f ON f.message_id = m.id
			WHERE m.conversation_id = ? AND m.status != ? AND mu.deleted_at IS NULL
			ORDER BY m.updated_at DESC`

	res := db.Raw(sql, model.TEXT,
		model.SYSTEM,
		model.FILE,
		userID,
		conversationID,
		model.RECALLED).
		Scan(&resp)

	if res.Error != nil {
		log.Println(res.Error)
		return nil, errors.New("服务器错误")
	}
	return resp, nil
}

// ConversationList 会话列表
func ConversationList(userID uint64) ([]model.ConversationListResp, error) {
	db := infra.GetDB()
	resp := make([]model.ConversationListResp, 0)

	sql := `SELECT cu.remark, 
       	cu.conversation_id,
       	cu.unread_count,
       	CASE m.status
		WHEN ? OR ? THEN t.text
		WHEN ? THEN f.file_name
		ELSE ''
		END AS content
		FROM conversation_users cu 
		LEFT JOIN messages m ON m.id = cu.last_message_id
		LEFT JOIN files f ON f.message_id = m.id
		LEFT JOIN texts t ON t.message_id = m.id
		WHERE cu.user_id = ? AND cu.deleted_at IS NULL
		ORDER BY cu.is_pinned DESC, cu.updated_at DESC `

	res := db.Raw(sql, model.TEXT,
		model.SYSTEM,
		model.FILE,
		userID).
		Scan(&resp)

	if res.Error != nil {
		log.Println(res.Error)
		return nil, errors.New("服务器错误")
	}

	return resp, nil
}

// getPrivateConversationID 获取两用户之间的私聊会话ID
// 两用户是好友关系才能正常工作，若不存在会话则创建新会话
func getPrivateConversationID(userID, friendID uint64) (uint64, error) {
	ok, err := isFriend(userID, friendID)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("两用户不是好友关系")
	}

	db := infra.GetDB()
	var cf model.ConversationFriend
	res := db.Model(&model.ConversationFriend{}).
		Select("conversation_id").
		Where("(user_id = ? AND friend_id = ?)"+
			"OR (friend_id = ? AND user_id = ?)", userID, friendID).
		First(&cf)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// 创建新会话
			return createPrivateConversation(userID, friendID)
		}
		log.Println(res.Error)
		return 0, errors.New("服务器错误")
	}

	return cf.ConversationID, nil
}

func createPrivateConversation(userID, friendID uint64) (uint64, error) {
	db := infra.GetDB()
	newID := utils.NewUniqueID()

	cf := model.ConversationFriend{
		ConversationID: newID,
		UserID:         userID,
		FriendID:       friendID,
	}

	res := db.Create(&cf)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return 0, errors.New("会话已存在")
		}
		log.Println(res.Error)
		return 0, errors.New("服务器错误")
	}

	return newID, nil
}

func CreateConversationUser(tx *gorm.DB, userID, conversationID uint64, remark string) error {
	res := tx.Model(&model.ConversationUser{}).
		Where("user_id = ? AND conversation_id = ?", userID, conversationID).
		Update("deleted_at", nil)
	if res.Error != nil {
		log.Println(res.Error)
		return errors.New("服务器错误")
	}
	if res.RowsAffected > 0 {
		// 已经存在该记录，直接返回
		return nil
	}

	cu := model.ConversationUser{
		UserID:         userID,
		ConversationID: conversationID,
		Remark:         remark,
	}

	res = tx.Create(&cu)
	if res.Error != nil {
		log.Println(res.Error)
		return errors.New("服务器错误")
	}

	return nil
}

func DeleteConversationUser(userID, conversationID uint64) error {
	db := infra.GetDB()
	res := db.Where("user_id = ? AND conversation_id = ?", userID, conversationID).
		Delete(&model.ConversationUser{})

	if res.Error != nil {
		log.Println(res.Error)
		return errors.New("服务器错误")
	}
	if res.RowsAffected == 0 {
		log.Println("删除 conversation_users 操作影响了0行表")
	}

	return nil
}
