package valid

import "testing"

func TestRule(t *testing.T) {
	r := NewRule().Set("Name,Age", Required, "eq=3", "le=1").Set("Age", "int", "test")
	sureMap := RM{"Name": "required,eq=3,le=1", "Age": "required,eq=3,le=1,int,test"}
	if !equal(r, sureMap) {
		t.Log(noEqErr)
	}
}

func TestRuleValid(t *testing.T) {
	type Tmp struct {
		Name      string
		Age       string
		ClassName string
	}
	v := Tmp{Name: "xue", Age: "12a"}
	validObj := NewVStruct()
	validObj.SetRule(NewRule().Set("Name,Age,ClassName", Required, "eq=3").Set("Age", "int"))
	if err := validObj.Valid(&v); err != nil {
		sureMsg := `"Tmp.Age" input "12a" is not integer; "Tmp.ClassName" input "" is required`
		if !equal(err.Error(), sureMsg) {
			t.Error(noEqErr)
		}
	}
}

func TestJoinTag2Val(t *testing.T) {
	if !equal(JoinTag2Val(Required, "", "必填"), Required+"|必填") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(Exist), Exist) {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(Either, "1"), Either+"=1") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(BothEq, "1"), BothEq+"=1") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VTo, "1~10", "需要在 1-10 的区间"), VTo+"=1~10|需要在 1-10 的区间") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VGe, "1", "大于或等于 1"), VGe+"=1|大于或等于 1") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VLe, "1"), VLe+"=1") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VOTo, "1~2"), VOTo+"=1~2") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VGt, "1", "大于 1"), VGt+"=1|大于 1") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VLt, "1", "小于 1"), VLt+"=1|小于 1") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VEq, "1"), VEq+"=1") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VNoEq, "1", "不等于 1"), VNoEq+"=1|不等于 1") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VIn, "1/2/3", "必须在 1,2,3 之中"), VIn+"=(1/2/3)|必须在 1,2,3 之中") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VInclude, "1/2", "包含 1,2"), VInclude+"=(1/2)|包含 1,2") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VPhone), VPhone) {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VEmail), VEmail) {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VIDCard, "", "非身份证"), VIDCard+"|非身份证") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VYear), VYear) {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VDate), VDate) {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VDatetime, "", "非日期时间"), VDatetime+"|非日期时间") {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VInt), VInt) {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VFloat), VFloat) {
		t.Error(noEqErr)
	}

	if !equal(JoinTag2Val(VRe, "[0-9]+", "必须为数字"), VRe+"='[0-9]+'|必须为数字") {
		t.Error(noEqErr)
	}
}
