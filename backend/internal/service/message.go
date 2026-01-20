package service

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"vvechat/internal/model"
	"vvechat/pkg/infra"
	"vvechat/pkg/utils"

	"gorm.io/gorm"
)

// sendMessageAuth 验证用户是否有权限在该会话中发送消息
func sendMessageAuth(userID, conversationID uint64) error {
	// 检查 conversation_users 表中是否存在该用户和会话
	var cnt int64
	db := infra.GetDB()
	err := db.Model(&model.ConversationUser{}).
		Where("user_id = ? AND conversation_id = ?", userID, conversationID).
		Count(&cnt).Error
	if err != nil {
		log.Println(err)
		return errors.New("服务器错误")
	}
	if cnt == 0 {
		return errors.New("无权限在该会话中发送消息")
	}
	return nil
}

// createSystemMessage 在一个会话中创建一个系统级消息
// newID用户指定该系统消息的ID
func createSystemMessage(tx *gorm.DB, content string, conversationID, newID uint64) error {
	newMsg := model.Message{
		SenderID:       0,
		ConversationID: conversationID,
		MyModel: model.MyModel{
			ID: newID,
		},
		Status: model.SYSTEM,
	}
	newText := model.Text{
		Text:      content,
		MessageID: newID,
	}
	res := tx.Create(&newMsg)
	if res.Error != nil {
		log.Println(res.Error)
		return errors.New("创建系统消息失败")
	}

	res = tx.Create(&newText)
	if res.Error != nil {
		log.Println(res.Error)
		return errors.New("创建系统消息失败")
	}

	return nil
}

// updateUnreadCount 给当前会话中除开当前sender的所有人的unread_count加一
func updateUnreadCount(tx *gorm.DB, senderID, conversationID uint64) error {
	res := tx.Model(&model.ConversationUser{}).
		Where("user_id != ? AND conversation_id = ?",
			senderID, conversationID).
		UpdateColumn("unread_count", gorm.Expr("unread_count + ?", 1))
	if res.Error != nil {
		log.Println(res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Println("更新unread count字段影响了0行表")
		return errors.New("更新unread count失败")
	}
	return nil
}

// updateLastMessageID 更新当前会话的last_message_id
func updateLastMessageID(tx *gorm.DB, conversationID, msgID uint64) error {
	res := tx.Model(&model.ConversationUser{}).
		Where("conversation_id = ?", conversationID).
		Update("last_message_id", msgID)
	if res.Error != nil {
		log.Println(res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Println("更新last msg id字段影响了0行表")
		return errors.New("更新last msg id失败")
	}
	return nil
}

func SendText(senderID, conversationID uint64, content string) (uint64, error) {
	err := sendMessageAuth(senderID, conversationID)
	if err != nil {
		return 0, err
	}

	newID := utils.NewUniqueID()
	newMsg := model.Message{
		SenderID:       senderID,
		ConversationID: conversationID,
		Status:         model.TEXT,
		MyModel: model.MyModel{
			ID: newID,
		},
	}
	newText := model.Text{
		Text:      content,
		MessageID: newID,
	}
	db := infra.GetDB()
	return newID, db.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&newMsg)
		if res.Error != nil {
			log.Println(res.Error)
			return errors.New("服务器错误")
		}

		res = tx.Create(&newText)
		if res.Error != nil {
			log.Println(res.Error)
			return errors.New("服务器错误")
		}

		err := updateLastMessageID(tx, conversationID, newID)
		if err != nil {
			return errors.New("服务器错误")
		}

		err = updateUnreadCount(tx, senderID, conversationID)
		if err != nil {
			return errors.New("服务器错误")
		}

		return nil
	})
}

func SendFile(senderID, conversationID uint64, file *multipart.FileHeader) (*model.SendFileResp, error) {
	err := sendMessageAuth(senderID, conversationID)
	if err != nil {
		return nil, err
	}

	newID := utils.NewUniqueID()

	// 获取 uploads 目录
	uploadDir := infra.GetFilePath()

	// 生成文件路径，使用newID作为文件名，保持原扩展名
	ext := filepath.Ext(file.Filename)
	fileName := strings.TrimSuffix(file.Filename, ext) // 原文件名
	filePath := filepath.Join(uploadDir, fmt.Sprintf("%d%s", newID, ext))

	// 保存文件
	if err := saveFile(file, filePath); err != nil {
		log.Println(err)
		return nil, errors.New("服务器错误")
	}

	// 获取文件信息
	fileSize := file.Size
	fileType := getFileType(filePath)

	newMsg := model.Message{
		SenderID:       senderID,
		ConversationID: conversationID,
		Status:         model.FILE,
		MyModel: model.MyModel{
			ID: newID,
		},
	}
	newFile := model.File{
		FileName:  fileName,
		FileType:  fileType,
		FileURL:   filePath,
		FileSize:  fileSize,
		MessageID: newID,
	}

	db := infra.GetDB()
	resp := &model.SendFileResp{
		MessageID: newID,
		FileName:  fileName,
		FileSize:  fileSize,
		FileType:  fileType,
	}
	return resp, db.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&newMsg)
		if res.Error != nil {
			log.Println(res.Error)
			return errors.New("服务器错误")
		}

		res = tx.Create(&newFile)
		if res.Error != nil {
			log.Println(res.Error)
			return errors.New("服务器错误")
		}

		err := updateLastMessageID(tx, conversationID, newID)
		if err != nil {
			return errors.New("服务器错误")
		}

		err = updateUnreadCount(tx, senderID, conversationID)
		if err != nil {
			return errors.New("服务器错误")
		}

		return nil
	})
}

func DownloadFile(userID, messageID uint64) (string, error) {
	db := infra.GetDB()

	// 一次查询完成：消息存在 + 用户在对话中 + 文件存在
	var file model.File
	err := db.Model(&model.File{}).
		Select("files.file_url").
		Joins("JOIN messages m ON m.id = files.message_id").
		Joins("JOIN conversation_users cu ON cu.conversation_id = m.conversation_id").
		Where("files.message_id = ? AND cu.user_id = ? AND m.status = ?",
			messageID, userID, model.FILE).
		First(&file).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("文件不存在或无访问权限")
		}
		log.Println("DB error:", err)
		return "", errors.New("服务器错误")
	}

	return file.FileURL, nil
}

func RecallMessage(userID, msgID uint64) (uint64, error) {
	db := infra.GetDB()
	var temp model.Message
	err := db.Model(&model.Message{}).
		Select("sender_id, conversation_id").
		Where("id = ?", msgID).
		First(&temp).
		Error
	if err != nil {
		log.Println(err)
		return 0, errors.New("服务器错误")
	}

	senderID := temp.SenderID
	conversationID := temp.ConversationID
	if senderID != userID {
		return 0, errors.New("不能撤回不是自己发的消息")
	}

	newID := utils.NewUniqueID()
	return newID, db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Message{}).
			Where("id = ?", msgID).
			Update("status", model.RECALLED)
		if res.Error != nil {
			log.Println(res.Error)
			return errors.New("服务器错误")
		}
		if res.RowsAffected == 0 {
			log.Println("撤回消息操作影响了0行表")
			return errors.New("服务器错误")
		}

		var senderName string
		err = tx.Model(&model.User{}).
			Where("id = ?", senderID).
			Pluck("name", &senderName).Error
		if err != nil {
			log.Println(err)
			return errors.New("服务器错误")
		}
		newContent := senderName + "撤回了一条消息"
		err = createSystemMessage(tx, newContent, conversationID, newID)
		if err != nil {
			return err
		}

		err = updateLastMessageID(tx, conversationID, newID)
		if err != nil {
			return err
		}
		return nil
	})
}

func DeleteMessage(userID, messageID uint64) error {
	db := infra.GetDB()
	var msg model.Message

	err := db.Model(&model.Message{}).
		Select("conversation_id").
		Where("id = ?", messageID).
		First(&msg).
		Error

	if err != nil {
		log.Println(err)
		return errors.New("服务器错误")
	}
	conversationID := msg.ConversationID

	return db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("user_id = ? AND message_id = ?", userID, messageID).
			Delete(&model.MessageUser{})

		if res.Error != nil {
			log.Println(res.Error)
			return errors.New("服务器错误")
		}
		if res.RowsAffected == 0 {
			log.Println("删除消息操作影响了0行表")
			return errors.New("服务器错误")
		}

		var lastID uint64
		sql := `SELECT m.id FROM messages m 
			LEFT JOIN message_users mu ON mu.message_id = m.id AND mu.user_id = ?
			WHERE m.status != ? AND mu.deleted_at IS NULL
			ORDER BY m.created_at DESC 
			LIMIT 1`
		res = tx.Raw(sql, userID, model.RECALLED).Scan(&lastID)
		if res.Error != nil {
			log.Println(res.Error)
			return errors.New("服务器错误")
		}
		if res.RowsAffected == 0 {
			log.Println("删除消息更新最后消息id 时没有查到id")
			return errors.New("服务器错误")
		}

		err = updateLastMessageID(tx, conversationID, lastID)
		if err != nil {
			return err
		}

		return nil
	})
}
