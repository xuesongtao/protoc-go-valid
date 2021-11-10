package file

import (
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
