package main

import (
	"fmt"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/pkg/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	initLogger()
	initEnv()
	initModule()
	web := InitWeb()
	web.Run(":9100")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func initModule() {
	if !service.App.Installed {
		return
	}
	webSettingSvc := InitWebSettingService()
	config, err := webSettingSvc.Read(service.App.Config)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("initModule#read app config fail"), zap.Error(err))
	}
	service.SetAppSetting(config)
	installSvc := InitInstallService()
	dao.SetGlobalDB(installSvc.CreateDB())
}

func initEnv() {
	service.App.Mode = service.DEV
	AppDir, err := WorkDir()
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	service.App.ConfDir = filepath.Join(AppDir, "/conf")
	service.App.Config = filepath.Join(service.App.ConfDir, "/app.ini")
	createDirIfNotExists(service.App.ConfDir)
	service.App.Installed = IsInstalled()
}

func WorkDir() (string, error) {
	if service.App.Mode == service.DEV {
		return utils.CurrentDir()
	}
	return utils.ExecDir()
}

func IsInstalled() bool {
	_, err := os.Stat(filepath.Join(service.App.ConfDir, "/install.lock"))
	if os.IsNotExist(err) {
		return false
	}
	return true
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
