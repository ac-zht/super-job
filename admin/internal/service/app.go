package service

import (
	"fmt"
	"github.com/ac-zht/super-job/admin/internal/repository"
	"github.com/ac-zht/super-job/admin/pkg/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

const (
	DEV  = "development"
	PROD = "production"
	TEST = "test"
)

func InitEnv() {
	repository.App.Mode = DEV
	AppDir, err := utils.WorkDir()
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	repository.App.ConfDir = filepath.Join(AppDir, "/conf")
	repository.App.Config = filepath.Join(repository.App.ConfDir, "/app.ini")
	createDirIfNotExists(repository.App.ConfDir)
	repository.App.Installed = IsInstalled()
}

func IsInstalled() bool {
	_, err := os.Stat(filepath.Join(repository.App.ConfDir, "/install.lock"))
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateInstallLock() error {
	_, err := os.Create(filepath.Join(repository.App.ConfDir, "/install.lock"))
	if err != nil {
		zap.L().Error("创建安装锁文件conf/install.lock失败")
	}
	return err
}

func createDirIfNotExists(path ...string) {
	for _, value := range path {
		if utils.FileExist(value) {
			continue
		}
		err := os.Mkdir(value, 0755)
		if err != nil {
			zap.L().Fatal(fmt.Sprintf("创建目录失败:%s", err.Error()))
		}
	}
}
