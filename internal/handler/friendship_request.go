package handler

import (
	"errors"
	"vvechat/internal/model"
	"vvechat/internal/service"
	"vvechat/pkg/judge"
	"vvechat/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 发送好友申请操作
func SendFriendRequest(c *gin.Context) {
	senderID := c.GetUint64("id")
	var req model.AddFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "json解析出错")
		return
	}

	err := service.SendFriendRequest(senderID, req.ReceiverID, req.VerificationMessage, req.SenderName)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			response.Fail(c, 400, "发送失败")
		} else if judge.IsUniqueConflict(err) {
			response.Fail(c, 409, "发送失败，请勿重复发送")
		} else {
			response.Fail(c, 500, "服务器出错")
		}
		return
	}

	response.Success(c, "发送成功", nil)
}

// 加载好友申请列表操作
func FriendRequestList(c *gin.Context) {
	receiverID := c.GetUint64("id")

	respSlice, err := service.FriendRequestList(receiverID)
	if err != nil {
		response.Fail(c, 500, "服务器错误")
		return
	}

	response.Success(c, "success", respSlice)
}
