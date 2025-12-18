package handler

import (
	"net/http"
	"vvechat/internal/model"
	"vvechat/internal/service"

	"vvechat/pkg/judge"
	"vvechat/pkg/response"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "json解析出错")
		return
	}

	user, err := model.NewUser(req.Name, req.Password, req.PhoneNumber)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	err = service.Register(user)
	if err != nil {
		if judge.IsUniqueConflict(err) {
			response.Fail(c, http.StatusBadRequest, "手机号已存在")
		} else {
			response.Fail(c, 500, "数据库错误")
		}
	} else {
		response.Success(c, "注册成功", nil)
	}
}

func LoginByUid(c *gin.Context) {
	var req model.LoginByUidReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "json解析出错")
		return
	}

	loginResp, err := service.LoginByUid(req.Uid, req.Password)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
	} else {
		response.Success(c, "登陆成功", loginResp)
	}
}

func LoginByPhone(c *gin.Context) {
	var req model.LoginByPhoneReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "json解析出错")
		return
	}

	loginResp, err := service.LoginByPhone(req.PhoneNumber, req.Password)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
	} else {
		response.Success(c, "登陆成功", loginResp)
	}
}
