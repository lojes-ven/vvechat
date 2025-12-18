package service

import (
	"errors"
	"vvechat/internal/model"
	"vvechat/pkg/infra"
	"vvechat/pkg/secure"

	"gorm.io/gorm"
)

func NewTokenResp(id uint64) (*model.TokenResp, error) {
	var resp model.TokenResp

	token, err := secure.NewToken(id)
	if err != nil {
		return nil, errors.New("token生成错误" + err.Error())
	}
	refreshToken, err := secure.NewRefreshToken(id)
	if err != nil {
		return nil, errors.New("refreshToken生成错误" + err.Error())
	}

	t := uint64(secure.GetExpiresTime().Seconds())
	if t <= 0 {
		return nil, errors.New("生成token时viper解析失败")
	}
	resp.ExpiresIn = t
	resp.Token, resp.RefreshToken = token, refreshToken

	return &resp, nil
}

func NewLoginResp(name string, uid string, id uint64) (*model.LoginResp, error) {
	var resp model.LoginResp

	tokenClass, err := NewTokenResp(id)
	if err != nil {
		return nil, err
	}

	resp.TokenClass.ExpiresIn = tokenClass.ExpiresIn
	resp.TokenClass.Token = tokenClass.Token
	resp.TokenClass.RefreshToken = tokenClass.RefreshToken

	resp.UserInfo.Name, resp.UserInfo.Uid = name, uid

	return &resp, nil
}

func GetUserByUid(uid string) (*model.User, error) {
	var user model.User
	res := infra.GetDB().Where("uid = ?", uid).First(&user)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("微信号不存在！")
		}
		return nil, res.Error
	}

	return &user, nil
}

func GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	res := infra.GetDB().Where("phone_number = ?", phone).First(&user)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("手机号不存在！")
		}
		return nil, res.Error
	}

	return &user, nil
}

// func IsUidExist(uid string) error {
// 	var cnt int64
// 	res := infra.GetDB().Model(&model.User{}).Where("uid = ?", uid).Count(&cnt)
// 	if res.Error != nil {
// 		return res.Error
// 	}

// 	exist := cnt > 0
// 	if exist {
// 		return gorm.ErrDuplicatedKey
// 	}
// 	return nil
// }

// func IsPhoneNumberExist(phone string) error {
// 	var cnt int64
// 	res := infra.GetDB().Model(&model.User{}).Where("phone_number = ?", phone).Count(&cnt)
// 	if res.Error != nil {
// 		return res.Error
// 	}

// 	exist := cnt > 0
// 	if exist {
// 		return gorm.ErrDuplicatedKey
// 	}
// 	return nil
// }

func Register(user *model.User) error {
	pwd, err := secure.HashString(user.Password)
	if err != nil {
		return err
	}

	user.Password = pwd

	return infra.GetDB().Create(user).Error
}

func LoginByUid(uid string, password string) (*model.LoginResp, error) {
	user, err := GetUserByUid(uid)
	if err != nil {
		return nil, errors.New("登陆失败 微信号或密码错误")
	}

	if ok := secure.VerifyPassword(user.Password, password); ok != nil {
		return nil, errors.New("登陆失败 微信号或密码错误")
	}

	return NewLoginResp(user.Name, user.Uid, user.ID)
}

func LoginByPhone(phone string, password string) (*model.LoginResp, error) {
	user, err := GetUserByPhone(phone)
	if err != nil {
		return nil, errors.New("登陆失败 手机号或密码错误")
	}

	if ok := secure.VerifyPassword(user.Password, password); ok != nil {
		return nil, errors.New("登陆失败 手机号或密码错误")
	}

	return NewLoginResp(user.Name, user.Uid, user.ID)
}
