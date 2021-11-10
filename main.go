package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/xuesongtao/protoc-go-valid/file"
	"gitee.com/xuesongtao/protoc-go-valid/log"
)

// handlePath 处理最后一个的路径服务
func handlePath(path string) string {
	lastSymbol := "/"
	if runtime.GOOS == "windows" {
		lastSymbol = "\\\\"
	}
	if strings.LastIndex(path, lastSymbol) != len(path)-1 {
		path += lastSymbol
	}
	return path
}

// handleDir 按目录处理
func handleDir(dirPath string) (isHasMatch bool) {
	dirs, err := os.ReadDir(dirPath)
	if err != nil {
		log.Error("os.ReadDir is failed, err: ", err)
		return
	}

	dirPath = handlePath(dirPath)
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		isHasMatch = true

		filename := dirPath + dir.Name()
		_ = handleFile(filename)
	}
	return
}

// handlePatternFiles 根据路径表达式处理
func handlePatternFiles(pattern string) (isHasMatch bool) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		log.Error("filepath.Glob is failed, err: ", err)
		return
	}

	for _, filename := range filenames {
		isHasMatch = true
		_ = handleFile(filename)
	}
	return
}

// handleFile 处理单个文件
func handleFile(filename string) (isHasMatch bool) {
	// 只处理 .go 文件
	if !strings.HasSuffix(filename, ".go") {
		return
	}
	isHasMatch = true

	log.Infof("parsing file %q for inject tag comments", filename)
	areas, err := file.ParseFile(filename)
	if err != nil {
		log.Error("file.ParseFile is failed, err: ", err)
		return
	}
	// log.Infof("areas: %+v", areas)

	if err = file.WriteFile(filename, areas); err != nil {
		log.Error("file.WriteFile is failed, err: ", err)
		return
	}
	log.Infof("file: %q is inject tag is success", filename)
	return
}

func main() {
	var inputDir, inputPattern, inputFile string
	flag.StringVar(&inputDir, "d", "", "注入的目录, 如: \"./proto\"")
	flag.StringVar(&inputPattern, "p", "", "注入匹配到的多个文件, 如: \"./*.pb.go\"")
	flag.StringVar(&inputFile, "f", "", "注入的单个文件, 如: \"xxx.pb.go\"")
	flag.Parse()

	var isHasMatch bool
	if inputDir != "" {
		isHasMatch = handleDir(inputDir)
	} else if inputPattern != "" {
		isHasMatch = handlePatternFiles(inputPattern)
	} else {
		isHasMatch = handleFile(inputFile)
	}

	if !isHasMatch {
		log.Error("it is not matched files, see: -help")
	}
}
