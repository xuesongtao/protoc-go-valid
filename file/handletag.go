package file

import (
	"fmt"
	"strings"
)

type tagItem struct {
	key   string
	value string
}

type tagItems []tagItem

func (t tagItems) format() string {
	tags := []string{}
	for _, item := range t {
		tags = append(tags, fmt.Sprintf(`%s:%s`, item.key, item.value))
	}
	return strings.Join(tags, " ")
}

// override 重写, 将输入的与现在的进行合并
func (t tagItems) override(inTags tagItems) tagItems {
	overridEd := []tagItem{}
	for i := range t {
		var dup = -1
		for j := range inTags {
			if t[i].key == inTags[j].key {
				dup = j
				break
			}
		}
		if dup == -1 {
			overridEd = append(overridEd, t[i])
		} else {
			overridEd = append(overridEd, inTags[dup])
			inTags = append(inTags[:dup], inTags[dup+1:]...)
		}
	}
	return append(overridEd, inTags...)
}

// newTagItems
func newTagItems(tag string) tagItems {
	items := []tagItem{}
	splitted := rTags.FindAllString(tag, -1)

	// fmt.Println("splitted: ", splitted)
	for _, t := range splitted {
		sepPos := strings.Index(t, ":")
		items = append(items, tagItem{
			key:   t[:sepPos],
			value: t[sepPos+1:],
		})
	}
	return items
}

// injectTag 注入 tag
func injectTag(contents []byte, area textArea) (injected []byte) {
	expr := make([]byte, area.End-area.Start)
	copy(expr, contents[area.Start-1:area.End-1])
	oldTag := newTagItems(area.CurrentTag)   // 原来的 tag
	injectTag := newTagItems(area.InjectTag) // 待注入的 tag
	finalTag := oldTag.override(injectTag)
	expr = rInject.ReplaceAll(expr, []byte(fmt.Sprintf("`%s`", finalTag.format())))
	injected = append(injected, contents[:area.Start-1]...)
	injected = append(injected, expr...)
	injected = append(injected, contents[area.End-1:]...)
	return
}

// tagFromComment 配注释中的注入 tag
func tagFromComment(comment string) (tag string) {
	match := rComment.FindStringSubmatch(comment)
	// fmt.Printf("comment: %s, match: %v\n", comment, match)
	if len(match) == 2 {
		tag = match[1]
	}
	return
}
