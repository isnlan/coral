package utils

import (
	"os"
	"strings"
)

// IsImg determines whether the specified extension is a image.
func IsImg(extension string) bool {
	ext := strings.ToLower(extension)

	switch ext {
	case ".jpg", ".jpeg", ".bmp", ".gif", ".png", ".svg", ".ico":
		return true
	default:
		return false
	}
}

func IsFile(path string) bool {
	fio, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	if nil != err {
		return false
	}

	return !fio.IsDir()
}

// 判断目录是否存在
func DirExists(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

//创建文件夹
func CreatedDir(dir string) error {
	exist := DirExists(dir)
	if !exist {
		// 创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

//移除文件夹
func RemoveDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

type DirCreator struct {
	dirs []string
}

func NewDirCreator() *DirCreator {
	return &DirCreator{dirs: make([]string, 0)}
}

func (dc *DirCreator) Push(dir string) {
	dc.dirs = append(dc.dirs, dir)
}

func (dc *DirCreator) Create() error {
	for _, dir := range dc.dirs {
		err := CreatedDir(dir)
		if err != nil {
			return err
		}
	}

	return nil
}
