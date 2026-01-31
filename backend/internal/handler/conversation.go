package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lojes7/inquire/internal/model"
	"github.com/lojes7/inquire/internal/service"
	"github.com/lojes7/inquire/pkg/response"
)

// StartPrivateConversation 新建私聊
func StartPrivateConversation(c *gin.Context) {
	userID := c.GetUint64("id")
	var req model.IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "json 解析出错")
		return
	}
	friendID := req.ID

	conversationID, err := service.StartPrivateConversation(userID, friendID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, 201, "success", conversationID)
}

// ChatHistoryList 加载聊天记录
func ChatHistoryList(c *gin.Context) {
	userID := c.GetUint64("id")
	id := c.Param("conversation_id")
	conversationID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Fail(c, 400, "conversation_id参数错误")
		return
	}

	resp, err := service.ChatHistoryList(userID, conversationID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, 200, "success", resp)
}

// ConversationList 会话列表
func ConversationList(c *gin.Context) {
	userID := c.GetUint64("id")

	resp, err := service.ConversationList(userID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, 200, "success", resp)
}
