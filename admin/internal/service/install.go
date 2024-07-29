package service

import (
	"context"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository"
)

type InstallService interface {
	Store(ctx context.Context, ins domain.Installation) error
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

func (svc *installService) Store(ctx context.Context, install domain.Installation) error {
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

func (svc *installService) pingDB(ins domain.Installation) error {
	var s domain.Setting
	s.DB.Host = ins.DbType
	s.DB.Port = ins.DbPort
	s.DB.User = ins.DbUsername
	s.DB.Password = ins.DbPassword
	s.DB.Database = ins.DbName
	s.DB.Charset = "utf8"

	return nil
}

func (svc *installService) Status(ctx context.Context) (bool, error) {
	return Installed, nil
}
