package main

import (
	"bytes"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/xuesongtao/protoc-go-valid/file"
	"gitee.com/xuesongtao/protoc-go-valid/log"
)

const (
	injectToolShellFileName  = "inject_tool.sh"
	injectToolSheellToolName = "protoc-go-valid-template" // 这里用于inject_tool.sh替换工具名
	windowsInjectTool        = "protoc-go-valid-windows"
	darwinInjectTool         = "protoc-go-valid-darwin"
	linuxInjectTool          = "protoc-go-valid-linux"
	tmpInjectTool            = "protoc-go-valid"
)

// copyInjectTool 将 inject_tool.sh 脚本移动到 GOPATH 下
func copyInjectTool() {
	goPath := os.Getenv("GOPATH")
	log.Info("GOPATH: ", goPath)
	if goPath == "" {
		log.Error("it is not found GOPATH, inject_tool.sh can not use")
		return
	}

	goBin := os.Getenv("GOBIN")
	log.Info("GOBIN: ", goBin)
	if goBin == "" {
		log.Error("it is not found GOBIN, inject_tool.sh can not use")
		return
	}

	// 对应操作系统
	toolSrc := ""
	switch runtime.GOOS {
	case "windows":
		toolSrc = windowsInjectTool
	case "darwin":
		toolSrc = darwinInjectTool
	default:
		toolSrc = linuxInjectTool
	}
	if err := file.CopyFile(toolSrc, goBin, true); err != nil {
		log.Errorf("copy %q to %q is failed, err: %v", toolSrc, goBin, err)
	}

	// 删除通用的
	_ = os.Remove(handlePath(goBin) + tmpInjectTool)

	// 判断下是否已经移动了, 如果已经移动就不处理了
	dest := handlePath(goPath) + injectToolShellFileName
	if _, err := os.Stat(dest); os.IsExist(err) {
		return
	}

	// 替换里面的工具名
	src := injectToolShellFileName
	if err := replaceFileContent(src, injectToolSheellToolName, toolSrc); err != nil {
		log.Error("replaceFileContent is failed, err: ", err)
		return
	}

	// 复制
	if err := file.CopyFile(src, dest); err != nil {
		log.Error("file.CopyFile is failed, err: ", err)
		return
	}
}

// replaceFileContent 文件内容替换
func replaceFileContent(src, old, new string) error {
	contentByte, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(src, bytes.ReplaceAll(contentByte, []byte(old), []byte(new)), fs.ModePerm)
}

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
	var (
		initProject                       bool
		inputDir, inputPattern, inputFile string
	)

	flag.BoolVar(&initProject, "init", false, "是否初始化项目, 如: protoc-go-valid -init=\"true\"")
	flag.StringVar(&inputDir, "d", "", "注入的目录, 如: protoc-go-valid -d \"./proto\"")
	flag.StringVar(&inputPattern, "p", "", "注入匹配到的多个文件, 如: protoc-go-valid -p \"./*.pb.go\"")
	flag.StringVar(&inputFile, "f", "", "注入的单个文件, 如: protoc-go-valid -f \"xxx.pb.go\"")
	flag.Parse()

	// 判断是否初始化
	if initProject {
		copyInjectTool()
		return
	}

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
