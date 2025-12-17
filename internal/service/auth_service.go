package service

import (
	"errors"
	"vvechat/internal/model"
)

func RefreshToken(typ string, id uint64) (*model.TokenResp, error) {
	if typ == "" || typ != "refresh" {
		return nil, errors.New("token格式有误")
	}

	return NewTokenResp(id)
}
