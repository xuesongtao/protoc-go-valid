package valid

import (
	"errors"
	"net/url"
	"reflect"
	"strings"
)

// VUrl 验证 url
type VUrl struct {
	ruleObj RM
	errBuf  *strings.Builder
	vc      *validCommon // 组合验证
}

// NewVUrl
func NewVUrl() *VUrl {
	obj := new(VUrl)
	obj.errBuf = newStrBuf()
	obj.vc = &validCommon{}
	return obj
}

// SetRule 添加验证规则
func (v *VUrl) SetRule(ruleObj RM) *VUrl {
	v.ruleObj = ruleObj
	return v
}

// Valid 验证
func (v *VUrl) Valid(src interface{}) error {
	if src == nil {
		return errors.New("src is nil")
	}

	var srcStr string
	switch v := src.(type) {
	case string:
		srcStr = v
	case *string:
		srcStr = *v
	default:
		return errors.New("src must is string/*string")
	}
	return v.validate(srcStr).getError()
}

// SetValidFn 自定义设置验证函数
func (v *VUrl) SetValidFn(validName string, fn CommonValidFn) *VUrl {
	v.vc.setValidFn(validName, fn)
	return v
}

// getValidFn 获取验证函数
func (v *VUrl) getValidFn(validName string) (CommonValidFn, error) {
	return v.vc.getValidFn(validName)
}

// validate 验证执行体
func (v *VUrl) validate(value string) *VUrl {
	// 解码处理
	decUrl, err := url.QueryUnescape(value)
	if err != nil {
		v.errBuf.WriteString(GetJoinFieldErr("", "", "url unescape is failed, err: "+err.Error()))
		return v
	}
	urlQuery := ""
	queryIndex := strings.Index(decUrl, "?")
	if queryIndex != -1 {
		urlQuery = decUrl[queryIndex+1:]
	}
	if urlQuery == "" {
		return v
	}

	var key, val string
	for _, query := range strings.Split(urlQuery, "&") {
		key = ""
		val = ""
		key2val := strings.Split(query, "=")
		l := len(key2val)
		if l > 0 {
			key = key2val[0]
		}
		if l > 1 {
			val = key2val[1]
		}

		validNames := v.ruleObj.Get(key)
		if validNames == "" {
			continue
		}
		// 根据验证内容进行验证
		for _, validName := range ValidNamesSplit(validNames) {
			if validName == "" {
				continue
			}

			validKey, _, cusMsg := ParseValidNameKV(validName)
			fn, err := v.getValidFn(validKey)
			if err != nil {
				v.errBuf.WriteString(GetJoinFieldErr("", key, err))
				continue
			}

			// 开始验证
			if fn == nil {
				switch validKey {
				case Required:
					if val != "" { // 验证必填
						continue
					}
					if cusMsg != "" {
						v.errBuf.WriteString(GetJoinValidErrStr("", key, "", cusMsg))
						continue
					}
					v.errBuf.WriteString(GetJoinValidErrStr("", key, "", ExplainEn, "it is", Required))
				case Either, BothEq:
					v.vc.initValid2FieldsMap(&name2Value{
						validName:  validName,
						fieldName:  key,
						cusMsg:     cusMsg,
						reflectVal: reflect.ValueOf(val),
					})
				default:
					v.errBuf.WriteString(GetJoinFieldErr("", key, "valid \""+validName+"\" is no support"))
				}
				continue

			}
			// 拓展的验证方法
			if val == "" { // 空就直接跳过
				continue
			}
			fn(v.errBuf, validName, "", key, reflect.ValueOf(val))
		}
	}
	return v
}

// getError 获取 err
func (v *VUrl) getError() error {
	defer putStrBuf(v.errBuf)
	v.vc.valid(v.errBuf)

	if v.errBuf.Len() == 0 {
		return nil
	}

	// 这里需要去掉最后一个 ErrEndFlag
	return errors.New(strings.TrimSuffix(v.errBuf.String(), ErrEndFlag))
}
