package service

import (
	"errors"
	"vvechat/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db}
}

func (s *UserService) GetUserByUid(uid string) (*model.User, error) {
	var user model.User
	res := s.db.Where("uid = ?", uid).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (s *UserService) IsUidExist(uid string) (bool, error) {
	var cnt int64
	res := s.db.Where("id = ?", uid).Count(&cnt)
	if res.Error != nil {
		return false, res.Error
	}

	if cnt > 0 {
		return true, nil
	}
	return false, nil
}

func (s *UserService) Register(user *model.User) error {
	if user.Uid == "" || user.Password == "" {
		return errors.New("密码或微信号不能为空！")
	}

	ok, err := s.IsUidExist(user.Uid)
	if err != nil {
		return err
	}
	if ok == true {
		return errors.New("微信号已存在！")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	return s.db.Create(user).Error
}
