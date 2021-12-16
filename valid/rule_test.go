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
	validObj.SetRule(RM{"name,age,classname": "required", "age": "int"})
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
	validObj.SetRule(NewRule().Set("name,age", "required").Set("age", "int"))
	t.Log(validObj.ruleMap.Get("name"))
	if err := validObj.Valid(&v); err != nil {
		t.Log(err)
	}
}
