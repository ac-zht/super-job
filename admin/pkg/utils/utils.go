package utils

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if os.IsPermission(err) {
		return false
	}
	return true
}

func WorkDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(execPath), nil
}

func RandAuthToken() string {
	buf := make([]byte, 32)
	_, err := crand.Read(buf)
	if err != nil {
		return RandString(64)
	}

	return fmt.Sprintf("%x", buf)
}

// RandString 生成长度为length的随机字符串
func RandString(length int64) string {
	sources := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sourceLength := len(sources)
	var i int64 = 0
	for ; i < length; i++ {
		result = append(result, sources[r.Intn(sourceLength)])
	}

	return string(result)
}
