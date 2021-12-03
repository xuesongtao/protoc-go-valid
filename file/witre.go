package file

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"gitee.com/xuesongtao/protoc-go-valid/log"
)

func WriteFile(inputPath string, areas []textArea) (err error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	// 处理 contents, 首先从文件的尾部注入自定义标记以保持顺序
	for i := 0; i < len(areas); i++ {
		area := areas[len(areas)-i-1]
		log.Infof("inject custom tag [%v] to expression [%v]", area.InjectTag, string(contents[area.Start-1:area.End-1]))
		contents = injectTag(contents, area)
	}

	if err = ioutil.WriteFile(inputPath, contents, 0644); err != nil {
		return
	}
	return
}

// CopyFile 复制文件
func CopyFile(src, dst string, isFirstDel ...bool) (err error) {
	if src == "" {
		return errors.New("source file cannot be empty")
	}

	if dst == "" {
		return errors.New("destination file cannot be empty")
	}

	// 如果相同就不处理
	if src == dst {
		return nil
	}

	// 删除原来的
	if len(isFirstDel) >0 && isFirstDel[0] {
		if err := os.Remove(dst); err != nil {
			return err
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() {
		if e := in.Close(); e != nil {
			err = e
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	// 复制
	if _, err = io.Copy(out, in); err != nil {
		return
	}

	// 写盘
	if err = out.Sync(); err != nil {
		return
	}

	// 调整权限
	if err = os.Chmod(dst, os.FileMode(0777)); err != nil {
		return
	}
	return
}
