package valid

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"
	"testing"
	"time"
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

func TestToString(t *testing.T) {
	t.Log(ToStr("1"))

	type Tmp struct {
		Name string
	}
	t.Log(ToStr(reflect.ValueOf(&Tmp{
		Name: "test",
	})))
}

func TestValidNamesSplit(t *testing.T) {
	t.Log(ValidNamesSplit("required|必填,phone|'手机号码必填,同时正确',re='\\d+{1,2}'"))
}

func TestLRU(t *testing.T) {
	var wg sync.WaitGroup
	size := 3
	lruSize := size - 1
	obj := NewLRU(lruSize)
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			time.Sleep(time.Second * time.Duration(num))
			obj.Store(num, "hello"+fmt.Sprint(num))
		}(i)
	}
	wg.Wait()
	t.Log(obj.Dump())
	if obj.Len() > lruSize {
		t.Error("size is not ok")
		return
	}

	if _, ok := obj.Load(0); ok {
		t.Error("0 is not ok")
		return
	}

	for i := size - lruSize; i < size; i++ {
		if v, _ := obj.Load(i); v.(string) != fmt.Sprintf("hello%d", i) {
			t.Errorf("%d is not ok", i)
			return
		}
	}
}
