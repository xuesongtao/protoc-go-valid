package valid



// *******************************************************************************
// *                              验证 struct                                     *
// *******************************************************************************

// Struct 验证结构体
func Struct(src interface{}, targetTag ...string) error {
	return NewVStruct(targetTag...).Valid(src)
}

// StructForFn 验证结构体, 同时设置自定义参数
func StructForFn(src interface{}, ruleObj RM, targetTag ...string) error {
	return NewVStruct(targetTag...).SetRule(ruleObj).Valid(src)
}

// ValidateStruct 验证结构体
func ValidateStruct(src interface{}, targetTag ...string) error {
	return NewVStruct(targetTag...).Valid(src)
}

// ValidStructForRule 自定义验证规则并验证
// 注: 通过字段名来匹配规则, 如果嵌套中如果有相同的名的都会走这个规则, 因此建议这种方式推荐使用非嵌套结构体
func ValidStructForRule(ruleObj RM, src interface{}, targetTag ...string) error {
	return NewVStruct(targetTag...).SetRule(ruleObj).Valid(src)
}

// ValidStructForMyValidFn 自定义单个验证函数
func ValidStructForMyValidFn(src interface{}, validName string, validFn CommonValidFn, targetTag ...string) error {
	return NewVStruct(targetTag...).SetValidFn(validName, validFn).Valid(src)
}




// *******************************************************************************
// *                             验证 单个变量                                     *
// *******************************************************************************

// Var 验证变量
func Var(src interface{}, rules ...string) error {
	return NewVVar().SetRules(rules...).Valid(src)
}

// VarForFn 验证变量, 同时设置自定义函数
func VarForFn(src interface{}, validFn CommonValidFn) error {
	return NewVVar().SetValidFn(validVarFieldName, validFn).Valid(src)
}



// *******************************************************************************
// *                             验证 query url                                   *
// *******************************************************************************

// Url 验证变量
func Url(src interface{}, ruleObj RM) error {
	return NewVUrl().SetRule(ruleObj).Valid(src)
}

// UrlForFn 验证 url, 同时设置自定义函数
func UrlForFn(src interface{}, validName string, validFn CommonValidFn) error {
	return NewVUrl().SetValidFn(validName, validFn).Valid(src)
}
