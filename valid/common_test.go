package valid

import (
	"regexp"
	"testing"
)

func TestParseValidNameKV(t *testing.T) {
	validName := "required|必填"
	k, v, m := ParseValidNameKV(validName)
	t.Logf("k: %q, v: %q, m: %q", k, v, m)
	if k != "required" || v != "" || m != "说明: 必填" {
		t.Error("parse is failed")
	}

	validName = "to=1~2|大于等于 1 且小于等于 2"
	k, v, m = ParseValidNameKV(validName)
	t.Logf("k: %q, v: %q, m: %q", k, v, m)
	if k != "to" || v != "1~2" || m != "说明: 大于等于 1 且小于等于 2" {
		t.Error("parse is failed")
	}
}

func TestRegexp(t *testing.T) {
	t.Log(regexp.MatchString("[\u4e00-\u9fa5]+", "abada"))
}
