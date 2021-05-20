package utils

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

// 将文件或目录打包成 .fileutil 文件
// src 是要打包的文件或目录的路径
// dstTar 是要生成的 .fileutil 文件的路径
// failIfExist 标记如果 dstTar 文件存在，是否放弃打包，如果否，则会覆盖已存在的文件
func Tar(src string, dstTar string, failIfExist bool) (err error) {
	// 清理路径字符串
	src = path.Clean(src)

	// 判断要打包的文件或目录是否存在
	if !Exists(src) {
		return errors.New("要打包的文件或目录不存在：" + src)
	}

	// 判断目标文件是否存在
	if FileExists(dstTar) {
		if failIfExist { // 不覆盖已存在的文件
			return errors.New("目标文件已经存在：" + dstTar)
		} else { // 覆盖已存在的文件
			if er := os.Remove(dstTar); er != nil {
				return er
			}
		}
	}

	// 创建空的目标文件
	fw, er := os.Create(dstTar)
	if er != nil {
		return er
	}
	defer fw.Close()

	// 创建 fileutil.Writer，执行打包操作
	tw := tar.NewWriter(fw)
	defer func() {
		// 这里要判断 tw 是否关闭成功，如果关闭失败，则 .fileutil 文件可能不完整
		if er := tw.Close(); er != nil {
			err = er
		}
	}()

	// 获取文件或目录信息
	fi, er := os.Stat(src)
	if er != nil {
		return er
	}

	// 获取要打包的文件或目录的所在位置和名称
	srcBase, srcRelative := path.Split(path.Clean(src))

	// 开始打包
	if fi.IsDir() {
		return tarDir(srcBase, srcRelative, tw, fi)
	} else {
		return tarFile(srcBase, srcRelative, tw, fi)
	}
}

// 因为要执行遍历操作，所以要单独创建一个函数
func tarDir(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	// 获取完整路径
	srcFull := srcBase + srcRelative

	// 在结尾添加 "/"
	last := len(srcRelative) - 1
	if srcRelative[last] != os.PathSeparator {
		srcRelative += string(os.PathSeparator)
	}

	// 获取 srcFull 下的文件或子目录列表
	fis, er := ioutil.ReadDir(srcFull)
	if er != nil {
		return er
	}

	// 开始遍历
	for _, fi := range fis {
		if fi.IsDir() {
			err := tarDir(srcBase, srcRelative+fi.Name(), tw, fi)
			if err != nil {
				return err
			}
		} else {
			err := tarFile(srcBase, srcRelative+fi.Name(), tw, fi)
			if err != nil {
				return err
			}
		}
	}

	// 写入目录信息
	if len(srcRelative) > 0 {
		hdr, er := tar.FileInfoHeader(fi, "")
		if er != nil {
			return er
		}
		hdr.Name = srcRelative

		if er = tw.WriteHeader(hdr); er != nil {
			return er
		}
	}

	return nil
}

// 因为要在 defer 中关闭文件，所以要单独创建一个函数
func tarFile(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	// 获取完整路径
	srcFull := srcBase + srcRelative

	// 写入文件信息
	hdr, er := tar.FileInfoHeader(fi, "")
	if er != nil {
		return er
	}
	hdr.Name = srcRelative

	if er = tw.WriteHeader(hdr); er != nil {
		return er
	}

	// 打开要打包的文件，准备读取
	fr, er := os.Open(srcFull)
	if er != nil {
		return er
	}
	defer fr.Close()

	// 将文件数据写入 tw 中
	if _, er = io.Copy(tw, fr); er != nil {
		return er
	}
	return nil
}

// 判断档案是否存在
func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// 判断文件是否存在
func FileExists(filename string) bool {
	fi, err := os.Stat(filename)
	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}
