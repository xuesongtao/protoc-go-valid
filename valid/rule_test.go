package valid

import "testing"

func TestRule(t *testing.T) {
	type Tmp struct {
		Name      string
		Age       string
		ClassName string
	}
	v := Tmp{Name: "xue", Age: "12a"}
	validObj := NewVStruct()
	validObj.SetRule(RM{"Name,Age,ClassName": "required", "Age": "int"})
	if err := validObj.Valid(&v); err != nil {
		t.Log(err)
	}
	
}

func TestRule2(t *testing.T) {
	type Tmp struct {
		Name string
		Age  string
	}
	v := Tmp{Name: "xue", Age: "12a"}
	validObj := NewVStruct()
	validObj.SetRule(NewRule().Set("Name,Age", Required, "eq=3").Set("Age", "int"))
	if err := validObj.Valid(&v); err != nil {
		t.Log(err)
	}
}

func TestRule3(t *testing.T) {
	r := NewRule().Set("Name,Age", Required, "eq=3", "le=1").Set("Age", "int", "test")
	t.Log(r)
}
