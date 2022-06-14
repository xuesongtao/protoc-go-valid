# protoc-go-valid [![Open Source Love](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](gitee.com/xuesongtao/protoc-go-valid) 

#### 1. 介绍

* 1. 通过对 `xxx.proto` 通过注释的形式加入验证 `tag`(使用方式文档下方有说明), 然后再使用 `inject_tool.sh xxx.proto` 编译, 这样生成的 `xxx.pb.go` 文件中的 `struct` 注入自定义的 `tag`
* 2. 通过验证器对 `struct` 中的 `tag` 进行验证

#### 2. 注入工具使用

* 1. 先下载本项目: `go get -u gitee.com/xuesongtao/protoc-go-valid`
* 2. `protoc-go-valid` 命令操作, 如下:   
    - 2.1 `protoc-go-valid -init="true"`
  + 2.1 `protoc-go-valid -d="待注入的目录"`
  + 2.2 `protoc-go-valid -p="匹配模式"`
  + 2.3 `protoc-go-valid -f="单个待注入的文件"`

* 3. 参考 `protoc-go-inject-tag`
	

#### 3. 工具补充

* 1. `protoc-go-valid -h` 可以通过这个查看帮助
* 2. 由于此操作是先执行 `protoc` 才再进行注入(需先安装 `protoc`), 项目中的 `inject_tool.sh` 整合了这两步操作, 可以执行 `protoc-go-valid -init="true"` 进行初始化操作, **说明:** 如果为 **windows** 需要使用 `powershell` 来执行, 如果失败的话, 可以直接将 `inject_tool.sh` 放到 GOPATH 下(主要是为了工具能命令行全局调用).
* 3. 根据自己的项目目录结构调整 `inject_tool.sh` 中 `proto` 和 `pb` 的目录, 相对于应用的目录; 如本项目, 修改如下下:

```
outPdProjectPath="test" # pb 放入的项目路径
protoFileDirName="test" # proto 存放的目录
```

#### 4. 验证器

##### 4.1 介绍

* 暂时只支持对 `struct` 的验证, **会一次性根据对 `struct` 设置的规则进行验证(包含嵌套验证), 将最终的所有错误都返回**
* 本工具易于拓展, 在功能方面不及业界著名 `validate`, 但在性能方面优于 `validate`, 需要调试可以在 `dev` 分支上, 同时该分支有性能测试数据

##### 4.2 验证 Struct

###### 4.2.1 支持的验证如下:  

| 标识     | 说明                                                                                                          |
| -------- | ------------------------------------------------------------------------------------------------------------ |
| required | 必填标识, 支持嵌套验证(结构体和结构体切片)                                                                         |
| exist    | 子对象有值才验证, 用于嵌套验证(结构体和结构体切片)                                                                   |
| either   | 多选一, 即多个中必须有一个必填, 格式为 "either=xxx"(通过数据进行标识)                                                 |
| botheq   | 多都相等, 即多个中必须都相等, 格式为 "botheq=xxx"(通过数据进行标识)                                                 |
| to       | 闭区间验证, 采用 `左右闭区间` , 格式为 "to=xxx\~xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度), 如: "to=1\~10"   |
| ge       | 大于或等于验证, 格式为 "ge=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度)                                     |
| le       | 小于或等于验证, 格式为: "le=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度)                                    |
| oto      | 开区间验证, 采用 `左右开区间` , 格式为 "oto=xxx\~xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度), 如: "oto=1\~10" |
| gt       | 大于验证, 格式为 "gt=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度)                                          |
| lt       | 小于验证, 格式为: "lt=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度)                                         |
| eq       | 等于验证, 格式为: "eq=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度)                                         |
| noeq     | 不等于验证, 格式为: "noeq=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度)                                     |
| in       | 指定输入选项, 格式为 "in=(xxx/xxx/xxx)", 如: "in=(1/abc/3)"                                                     |
| include  | 指定输入包含选项, 格式为 "include=(xxx/xxx/xxx)", 如: "in=(hello/2/3)"                                           |
| phone    | 手机号验证                                                                                                     |
| email    | 邮箱验证                                                                                                       |
| idcard   | 身份证号码验证                                                                                                  |
| year     | 年验证                                                                                                         |
| year2month| 年月验证, 支持分割符, 默认按照"-". 验证:xxxx/xx, 格式: "year2month=/"                                             |
| date     | 日期验证, 支持分割符, 默认按照"-". 验证:xxxx/xx/xx, 格式: "date=/"                                                 |
| datetime | 时间验证, 支持分割符, 默认按照"-". 验证:xxxx/xx/xx xx:xx:xx, 格式: "datetime=/"                                     |
| int      | 整数型验证(字段类型为字符串)                                                                                      |
| float    | 浮动数型验证(字段类型为字符串)                                                                                    |

###### 4.2.2 设置验证

* 1. 通过设置 `tag` 进行设置验证规则, 默认目标为 `valid`
* 2. 支持通过创建 `RM` 对象进行自定义设置验证规则, 其验证优先级高于 `xxx.pb.go` 里的规则, `RM` 暂不支持嵌套


###### 4.2.3 其他

* 1. 默认按照 `tag` 进行处理, 如果设置 `RM` 对象会以此规则为准
* 2. 如果验证方法没有实现的, 可以调用 `SetCustomerValidFn` 自定义
* 3. 使用的可以参考 `example_test.go` 和 `valid_test.go`


#### 5 使用示例:

* `proto` 内容如下: 

```
message Man {
    string name = 1; // 姓名 @tag valid:"required,to=1~3" 
    int32 age = 2; // 年龄 @tag valid:"to=1~150"
}
```

* **注:** 编写 `xxx.proto` 时, 需要加将 `@tag xxx` 放到注释的最后面

* 执行命令: `inject_tool.sh xxx.proto` 生成 `pd` 内容如下: 

```
type Man struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" valid:"required,to=1~3"` // 姓名 @tag valid:"required,to=1~3"
	Age  int32  `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty" valid:"to=1~150"`  // 年龄 @tag valid:"to=1~150"
}
```

* 代码里的使用

```
	m := &test.Man{
			Name: "xue",
			Age:  -1,
	}
	fmt.Println(ValidateStruct(m))

	// Output: "Man.Age" input "-1" is size less than 1
```

#### 最后

* 欢迎大佬们指正, 同时也希望大佬给 **star**
