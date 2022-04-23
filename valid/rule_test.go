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
