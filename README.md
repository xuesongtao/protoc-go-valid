# proto 中注入 tag
#### 1. 介绍
- 1. 对 `xxx.pd.go` 文件中的 `struct` 注入自定义的 `tag`
- 2. 通过验证器对内容进行验证, 验证器暂只支持: 必填, 长度, 多个一个必填


#### 2. 注入工具使用
- 1. 先下载本项目: `go get -u gitee.com/xuesongtao/protoc-go-valid`
- 2. `protoc-go-valid` 命令操作, 如下: 
    - 2.1 `protoc-go-valid -init="true"`
	- 2.1 `protoc-go-valid -d="待注入的目录"`
	- 2.2 `protoc-go-valid -p="匹配模式"`
	- 2.3 `protoc-go-valid -f="单个待注入的文件"`
- 3. 参考 `protoc-go-inject-tag`
	

#### 3. 工具补充
- 1. protoc-go-valid -h` 可以通过这个查看帮助
- 2. 由于此操作是先执行 `protoc` 才再进行注入(需先按照 `protoc`), 项目中的 `inject_tool.sh` 整合了这两步操作, 可以执行 `protoc-go-valid -init="true"`, **说明:** 如果为 **windows** 需要使用 `powershell` 来执行
- 3. 根据直接项目目录结构调整 `inject_tool.sh` 中 `proto` 和 `pb` 的目录, 相对于应用的目录; 如本项目, 修改如下下:
```
outPdProjectPath="test" # pb 放入的项目路径
protoFileDirName="test" # proto 存放的目录
```


#### 4. 验证器
##### 4.1 验证器支持的验证如下:
  
| 标识     | 说明                                                                                             |
| -------- | ------------------------------------------------------------------------------------------------ |
| required | 必填标识, 支持嵌套验证(结构体和结构体切片)                                                       |
| either   | 多选一, 即多个中必须有一个必填, 格式为 either=xxx(通过数据进行标识)                              |
| to       | 长度验证, 格式为 "to=xxx~xxx"(字段类型: 字符串则为长度, 字段类型: 数字型则为大小), 如: "to=1~10" |
| in       | 指定输入选项, 格式为 in=(xxx/xxx/xxx), 如: in(1/2/3)                                             |
| phone    | 手机号验证                                                                                       |
| email    | 邮箱验证                                                                                         |
| idcard   | 身份证号码验证                                                                                   |


##### 4.2 使用示例:
-  `pd` 内容如下: 
```
type Man struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" valid:"required"` // 姓名 @tag valid:"required"
	Age  int32  `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty" valid:"to=1~150"`  // 年龄 @tag valid:"to=1~150"
}
```
- **注:** 编写 `xxx.proto` 时, 需要加将 `@tag xxx` 放到注释的最后面

- 代码里的使用
```
	m := &test.Man{
			Name: "xue",
			Age:  0,
	}
	t.Log(ValidateStruct(m))

	// 输出: "Man.Age" is size less than 1
```


