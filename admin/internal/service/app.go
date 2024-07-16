package service

import (
	"fmt"
	"github.com/zc-zht/super-job/admin/pkg/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

var (
	ConfDir   string
	Installed bool
)

func InitEnv() {
	AppDir, err := utils.WorkDir()
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	ConfDir = filepath.Join(AppDir, "/conf")
	createDirIfNotExists(ConfDir)
	Installed = IsInstalled()
}

func IsInstalled() bool {
	_, err := os.Stat(filepath.Join(ConfDir, "/install.lock"))
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateInstallLock() error {
	_, err := os.Create(filepath.Join(ConfDir, "/install.lock"))
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
