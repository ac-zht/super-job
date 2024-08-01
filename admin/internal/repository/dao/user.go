package dao

import (
	"context"
	"time"
)

type UserDAO interface {
	List(ctx context.Context, offset, limit int) ([]User, error)
	GetById(ctx context.Context, id int64) (User, error)
	Insert(ctx context.Context, u User) (int64, error)
	Update(ctx context.Context, u User) error
	Delete(ctx context.Context, id int64) error
}

type GORMUserDAO struct {
	BaseModel
}

func NewUserDAO(base BaseModel) UserDAO {
	return &GORMUserDAO{
		BaseModel: base,
	}
}

func (dao *GORMUserDAO) List(ctx context.Context, offset, limit int) ([]User, error) {
	var users []User
	err := dao.DB().WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	return users, err
}

func (dao *GORMUserDAO) GetById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.DB().WithContext(ctx).First(&u, id).Error
	return u, err
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) (int64, error) {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.DB().WithContext(ctx).Create(&u).Error
	return u.Id, err
}

func (dao *GORMUserDAO) Update(ctx context.Context, u User) error {
	u.Utime = time.Now().UnixMilli()
	return dao.DB().WithContext(ctx).Updates(&u).Error
}

func (dao *GORMUserDAO) Delete(ctx context.Context, id int64) error {
	return dao.DB().WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Name     string `gorm:"type:varchar(32);unique"`
	Password string `gorm:"type:char(32)"`
	Email    string `gorm:"type:varchar(50);unique"`
	Salt     string `gorm:"char(6)"`
	IsAdmin  uint8  `gorm:"type:tinyint"`
	Status   uint8  `gorm:"type:tinyint"`
	Ctime    int64
	Utime    int64
}

const (
	Disabled uint8 = iota
	Enabled
)

const PasswordSaltLength = 6
