package file

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"runtime"
	"strings"
)

const (
	InjectTagFlag = "@tag" // 注入 tag 的标识
)

var (
	rComment = regexp.MustCompile(`@tag (.*)`) // 匹配注入 tag
	rInject  = regexp.MustCompile("`.+`$")
	rTags    = regexp.MustCompile(`\w+:"[^"]+"`) // 匹配 tag
)

// textArea
type textArea struct {
	Start      int    // 开始位置
	End        int    // 截止位置
	CurrentTag string // 已有 tag
	InjectTag  string // 注入的 tag
}

// HandlePath 处理最后一个的路径服务
func HandlePath(path string) string {
	lastSymbol := "/"
	if runtime.GOOS == "windows" {
		lastSymbol = "\\\\"
	}
	if strings.LastIndex(path, lastSymbol) != len(path)-1 {
		path += lastSymbol
	}
	return path
}

// ParseFile 解析文件
func ParseFile(inputPath string) (areas []textArea, err error) {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, inputPath, nil, parser.ParseComments)
	if err != nil {
		return
	}
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var typeSpec *ast.TypeSpec // 类型
		for _, spec := range genDecl.Specs {
			if ts, tsOK := spec.(*ast.TypeSpec); tsOK {
				typeSpec = ts
				break
			}
		}

		// 空就跳过
		if typeSpec == nil {
			continue
		}

		// 不是结构体就跳过
		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		for _, field := range structDecl.Fields.List {
			var comments []*ast.Comment
			// 字段的注释
			if field.Comment != nil {
				comments = append(comments, field.Comment.List...)
			}

			// 组装数据
			for _, comment := range comments {
				tag := tagFromComment(comment.Text)
				if tag == "" {
					continue
				}

				currentTag := field.Tag.Value
				area := textArea{
					Start:      int(field.Pos()),
					End:        int(field.End()),
					CurrentTag: currentTag[1 : len(currentTag)-1], // 去掉 ``
					InjectTag:  tag,
				}
				areas = append(areas, area)
			}
		}
	}
	return
}
