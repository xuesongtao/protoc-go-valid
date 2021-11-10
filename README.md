# proto 中注入 tag
#### 1. 介绍
- 1.对 `xxx.proto` 文件注入 tag 
- 2.通过验证器对内容进行验证
- **注:** 必须先按照 `protoc`


#### 2. 注入工具使用
- 1.先下载本项目: `go get -u gitee.com/xuesongtao/gitee.com/xuesongtao/protoc-go-valid`
- 2.把项目中的`cjpro.sh`加到环境变量, **说明:** 如果为 **windows** 需要使用 `powershell` 来执行


#### 3. 验证器使用
- 1.请先下载项目
- `pd` 内容如下: 
```
type Man struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" valid:"required"` // 姓名 @tag valid:"required"
	Age  int32  `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty" valid:"to=1~150"`  // 年龄 @tag valid:"to=1~150"
}
```

- 代码里的使用
```
	m := &test.Man{
			Name: "xue",
			Age:  0,
	}
	t.Log(ValidateStruct(m))

	// 输出: "Man.Age" is size less than 1
```