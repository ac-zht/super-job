package domain

type Executor struct {
	Id    int64
	Name  string
	Hosts []string
	Ctime int64
	Utime int64
}
