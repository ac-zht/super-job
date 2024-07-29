package repository

import "gorm.io/gorm"

type InstallRepository interface {
	CreateTmpDB() (*gorm.DB, error)
	CreateDB() error
}

type installRepository struct {
}

func (i installRepository) CreateTmpDB() (*gorm.DB, error) {
	//TODO implement me
	panic("implement me")
}

func (i installRepository) CreateDB() error {
	//TODO implement me
	panic("implement me")
}
