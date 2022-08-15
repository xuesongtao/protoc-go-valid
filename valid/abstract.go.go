package valid

type Valider interface {
	SetRule(RM) Valider                       // 设置规则
	SetValidFn(string, CommonValidFn) Valider // 设置自定义参数
	Valid(interface{}) error                  // 验证
}

