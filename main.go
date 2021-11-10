package main

import (
	"flag"
	"os"
	"protoc-go-cjvalid/file"
	"protoc-go-cjvalid/log"
	"runtime"
	"strings"
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

func main() {
	var inputFiles string
	flag.StringVar(&inputFiles, "f", "", "pattern to match input file(s)")
	flag.Parse()

	if len(inputFiles) == 0 {
		log.Error("f file is mandatory, see: -help")
		return
	}

	dirs, err := os.ReadDir(inputFiles)
	if err != nil {
		log.Error(err)
		return
	}

	inputFiles = handlePath(inputFiles)
	var matched int
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}

		// 只处理 .go 文件
		if !strings.HasSuffix(dir.Name(), ".go") {
			continue
		}
		matched++

		filename := inputFiles + dir.Name()
		log.Infof("parsing file %q for inject tag comments", filename)
		areas, err := file.ParseFile(filename)
		if err != nil {
			log.Error(err)
			return
		}
		// log.Infof("areas: %+v", areas)

		if err = file.WriteFile(filename, areas); err != nil {
			log.Error(err)
			return
		}
		log.Infof("%q is inject tag is success", filename)
	}

	if matched == 0 {
		log.Error("f %q matched no files, see: -help", inputFiles)
	}
}
