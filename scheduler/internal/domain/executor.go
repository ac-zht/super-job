package domain

import (
	"github.com/ac-zht/super-job/scheduler/internal/repository/dao"
	"strings"
)

type Executor struct {
	Id    int64
	Name  string
	Hosts []string
}

func ToDomain(e dao.Executor) Executor {
	return Executor{
		Id:    e.Id,
		Name:  e.Name,
		Hosts: strings.Split(e.Hosts, ","),
	}
}

func ToEntity(e Executor) dao.Executor {
	return dao.Executor{
		Id:    e.Id,
		Name:  e.Name,
		Hosts: strings.Join(e.Hosts, ","),
	}
}
