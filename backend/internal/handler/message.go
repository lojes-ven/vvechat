package handler

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lojes7/inquire/internal/model"
	"github.com/lojes7/inquire/internal/service"
	"github.com/lojes7/inquire/pkg/infra"
	"github.com/lojes7/inquire/pkg/response"
)

func SendText(c *gin.Context) {
	senderID := c.GetUint64("id")
	var req model.SendTextReq
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, 400, "json 解析出错")
		return
	}
	conversationID := req.ConversationID
	content := req.Content

	msgID, err := service.SendText(senderID, conversationID, content)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, 201, "success", msgID)
}

func SendFile(c *gin.Context) {
	userID := c.GetUint64("id")

	conversationIDStr := c.PostForm("conversation_id")
	if conversationIDStr == "" {
		response.Fail(c, 400, "conversation_id 是空的")
		return
	}
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "conversation_id 格式错误")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 400, "没有接收到文件")
		return
	}

	resp, err := service.SendFile(userID, conversationID, file)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, 201, "success", resp)
}

func DownloadFile(c *gin.Context) {
	userID := c.GetUint64("id")
	messageIDStr := c.Param("message_id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)
	if err != nil {
		response.Fail(c, 400, "message_id 格式错误")
		return
	}

	fileURL, err := service.DownloadFile(userID, messageID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	// 安全检查，确保文件路径在允许的目录下
	allowedDir := infra.GetFilePath()
	if !strings.HasPrefix(fileURL, allowedDir) {
		response.Fail(c, 403, "非法文件路径")
		return
	}

	// 打开文件
	file, err := os.Open(fileURL)
	if err != nil {
		if os.IsNotExist(err) {
			response.Fail(c, 404, "文件不存在")
		} else {
			log.Printf("打开文件失败: %v\n", err)
			response.Fail(c, 500, "文件读取失败")
		}
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		response.Fail(c, 500, "无法获取文件信息")
		return
	}

	// 设置响应头
	fileName := filepath.Base(fileURL) // 例如 "report.pdf"

	// 自动推断 MIME 类型
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	c.Header("Content-Type", mimeType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.PathEscape(fileName)))

	// 流式传输文件（不加载全文件到内存）
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		// 客户端可能取消了下载，通常不用报错
		log.Printf("文件传输中断: %v", err)
		return
	}
}

func RecallMessage(c *gin.Context) {
	userID := c.GetUint64("id")
	var req model.IDReq
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, 400, "json 解析错误")
		return
	}
	// 撤回消息会创建一个系统级消息，这里拿到该消息的ID，返回给前端
	systemMsgID, err := service.RecallMessage(userID, req.ID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, 201, "success", systemMsgID)
}

func DeleteMessage(c *gin.Context) {
	userID := c.GetUint64("id")
	var req model.IDReq
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, 400, "json 解析错误")
		return
	}

	err := service.DeleteMessage(userID, req.ID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, 200, "success", nil)
}
