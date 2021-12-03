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
	injectToolTemplateSh = "inject_tool_template.sh"
	injectToolSh         = "inject_tool.sh"
	injectToolReplaceTag = "protoc-go-valid-template" // 这里用于inject_tool_template.sh替换工具名
	windowsInjectTool    = "protoc-go-valid-windows"
	darwinInjectTool     = "protoc-go-valid-darwin"
	linuxInjectTool      = "protoc-go-valid-linux"
	tmpInjectTool        = "protoc-go-valid"
)

// copyInjectTool 将 inject_tool.sh 脚本移动到 GOPATH 下
func copyInjectTool() {
	goPath := os.Getenv("GOPATH")
	log.Info("GOPATH: ", goPath)
	if goPath == "" {
		log.Error("it is not found GOPATH, inject_tool.sh can not use")
		return
	}

	goBin := file.HandlePath(goPath) + "bin"

	// 将操作系统对应的注入工具复制到 GOBIN 下
	toolSrc := ""
	switch runtime.GOOS {
	case "windows":
		toolSrc = windowsInjectTool
	case "darwin":
		toolSrc = darwinInjectTool
	default:
		toolSrc = linuxInjectTool
	}
	if err := file.CopyFile(toolSrc, file.HandlePath(goBin)+toolSrc, true); err != nil {
		log.Errorf("copy %q to %q is failed, err: %v", toolSrc, goBin, err)
		return
	}

	// 删除通过 go install 生成的可执行文件
	_ = os.Remove(file.HandlePath(goBin) + tmpInjectTool)

	// 判断下是否已经移动了, 如果已经移动就不处理了
	dest := file.HandlePath(goPath) + injectToolSh
	if _, err := os.Stat(dest); os.IsExist(err) {
		return
	}

	// 替换里面的工具名
	src := injectToolSh
	if err := createInjectToolSh(injectToolTemplateSh, src, injectToolReplaceTag, toolSrc); err != nil {
		log.Error("createInjectToolSh is failed, err: ", err)
		return
	}

	// 复制
	if err := file.CopyFile(src, dest); err != nil {
		log.Error("file.CopyFile is failed, err: ", err)
		return
	}
}

// createInjectToolSh 文件内容替换
func createInjectToolSh(template, create, old, new string) error {
	contentByte, err := os.ReadFile(template)
	if err != nil {
		return err
	}
	return os.WriteFile(create, bytes.ReplaceAll(contentByte, []byte(old), []byte(new)), fs.ModePerm)
}

// handleDir 按目录处理
func handleDir(dirPath string) (isHasMatch bool) {
	dirs, err := os.ReadDir(dirPath)
	if err != nil {
		log.Error("os.ReadDir is failed, err: ", err)
		return
	}

	dirPath = file.HandlePath(dirPath)
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
