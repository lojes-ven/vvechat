package handler

import (
	"errors"
	"log"
	"strconv"
	"vvechat/internal/model"
	"vvechat/internal/service"
	"vvechat/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FriendshipList(c *gin.Context) {
	id := c.GetUint64("id")
	resp, err := service.FriendshipList(id)

	if err != nil {
		response.Fail(c, 500, "服务器错误")
		return
	}

	response.Success(c, 200, "success", resp)
}

func DeleteFriendship(c *gin.Context) {
	userID := c.GetUint64("id")
	friendId := c.Param("friend_id")

	friendID, err := strconv.ParseUint(friendId, 10, 64)
	if err != nil {
		response.Fail(c, 400, "friend_id不合法")
		return
	}

	err = service.DeleteFriendship(userID, friendID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, 400, "好友不存在")
		} else {
			response.Fail(c, 500, "服务器错误")
		}
		return
	}

	response.Success(c, 201, "success", nil)
}

func ReviseRemark(c *gin.Context) {
	userID := c.GetUint64("id")
	id := c.Param("friend_id")
	friendID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("修改备注时，好友id转为uint64类型时错误")
		response.Fail(c, 400, "修改失败")
		return
	}

	var req model.RemarkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("json解析错误")
		response.Fail(c, 400, "输入不合法")
		return
	}

	err = service.ReviseRemark(userID, friendID, req.Remark)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.Success(c, 200, "success", nil)
}
