package repository

import (
	"errors"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"gopkg.in/ini.v1"
)

type WebSettingRepository interface {
	Read(fn string) (*domain.Setting, error)
	Write(config []string, fn string) error
}

type webSettingRepository struct {
}

func NewWebSettingRepository() WebSettingRepository {
	return &webSettingRepository{}
}

func (w *webSettingRepository) Read(fn string) (*domain.Setting, error) {
	config, err := ini.Load(fn)
	if err != nil {
		return nil, err
	}
	section := config.Section(domain.DefaultSection)
	var s domain.Setting
	s.DB.Engine = section.Key("db.engine").MustString("mysql")
	s.DB.Host = section.Key("db.host").MustString("127.0.0.1")
	s.DB.Port = section.Key("db.port").MustInt(3306)
	s.DB.User = section.Key("db.user").MustString("")
	s.DB.Password = section.Key("db.password").MustString("")
	s.DB.Database = section.Key("db.database").MustString("super-job")
	s.DB.Prefix = section.Key("db.prefix").MustString("")
	s.DB.Charset = section.Key("db.charset").MustString("utf8")
	s.DB.MaxIdleConns = section.Key("db.max.idle.conns").MustInt(30)
	s.DB.MaxOpenConns = section.Key("db.max.open.conns").MustInt(100)
	return &s, nil
}

func (w *webSettingRepository) Write(config []string, fn string) error {
	if len(config) == 0 {
		return errors.New("params is empty")
	}
	if len(config)%2 != 0 {
		return errors.New("param mismatching")
	}
	file := ini.Empty()
	section, err := file.NewSection(domain.DefaultSection)
	if err != nil {
		return err
	}
	for i := 0; i < len(config); {
		_, err = section.NewKey(config[i], config[i+1])
		if err != nil {
			return err
		}
		i += 2
	}
	return file.SaveTo(fn)
}
