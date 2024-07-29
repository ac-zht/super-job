package service

import (
	"context"
	"github.com/ac-zht/super-job/admin/internal/repository"
)

type InstallService interface {
	Store(ctx context.Context) error
	Status(ctx context.Context) (bool, error)
}

type installService struct {
	settingRepo repository.SettingRepository
}

func NewInstallService(setRepo repository.SettingRepository) InstallService {
	return &installService{
		settingRepo: setRepo,
	}
}

func (svc *installService) Store(ctx context.Context) error {
	//ping数据库
	//根据提交的消息写数据库配置文件
	//读取文件到内存
	//创建数据库
	//生成表
	//初始化配置表字段
	err := svc.settingRepo.InitBasicField(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (svc *installService) pingDb() error {
	return nil
}

func (svc *installService) Status(ctx context.Context) (bool, error) {
	return Installed, nil
}
