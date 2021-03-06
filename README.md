# [protoc-go-valid](https://gitee.com/xuesongtao/protoc-go-valid) 
[![OSCS Status](https://www.oscs1024.com/platform/badge/xuesongtao/protoc-go-valid.svg?size=small)](https://www.oscs1024.com/project/xuesongtao/protoc-go-valid?ref=badge_small)

#### 项目背景

*  1. 在 protobuf 方面验证器常用的为 `go-proto-validators` 验证器, 使用方面个人认为较为繁琐，代码量比较多, 使用如下:  
 

```
syntax = "proto3";
package validator.examples;
import "github.com/mwitkow/go-proto-validators/validator.proto";

message InnerMessage {
    // some_integer can only be in range (0, 100).
    int32 some_integer = 1 [(validator.field) = {int_gt: 0, int_lt: 100}];
    // some_float can only be in range (0;1).
    double some_float = 2 [(validator.field) = {float_gte: 0, float_lte: 1}];
}
```

*  2. 本验证器, 相同功能使用如下:  

```
syntax = "proto3";
package examples;

message InnerMessage {
    // some_integer can only be in range (0, 100).
    int32 some_integer = 1; // @tag oto=0~100
    // some_float can only be in range (0;1).
    double some_float = 2; // @tag oto=0~1
}
```

#### 1. 介绍

*  1. 通过对 `xxx.proto` 通过注释的形式加入验证 `tag`(使用方式文档下方有说明), 然后再使用 `inject_tool.sh xxx.proto` 编译, 这样生成的 `xxx.pb.go` 文件中的 `struct` 注入自定义的 `tag`

*  2. 通过验证器对 `struct` 中的 `tag` 进行验证

#### 2. 注入工具使用

*  1. 先下载本项目: `go get -u gitee.com/xuesongtao/protoc-go-valid`

*  2.  `protoc-go-valid` 命令操作, 如下:  

  + 2.1 `protoc-go-valid -init="true"`
  + 2.2 `protoc-go-valid -d="待注入的目录"`
  + 2.3 `protoc-go-valid -p="匹配模式"`
  + 2.4 `protoc-go-valid -f="单个待注入的文件"`

*  3. 参考 `protoc-go-inject-tag`

#### 3. 工具补充

*  1.  `protoc-go-valid -h` 可以通过这个查看帮助

*  2. 由于此操作是先执行 `protoc` 才再进行注入(需先安装 `protoc`), 项目中的 `inject_tool.sh` 整合了这两步操作, 可以执行 `protoc-go-valid -init="true"` 进行初始化操作, **说明:** 如果为 **windows** 需要使用 `powershell` 来执行, 如果失败的话, 可以直接将 `inject_tool.sh` 放到 GOPATH 下(主要是为了工具能命令行全局调用).

*  3. 根据自己的项目目录结构调整 `inject_tool.sh` 中 `proto` 和 `pb` 的目录, 相对于应用的目录; 如本项目, 修改如下下:  

```
outPdProjectPath="test" # pb 放入的项目路径
protoFileDirName="test" # proto 存放的目录
```

#### 4. 验证器

##### 4.1 介绍

* 支持对 **一个或多个struct**, **会一次性根据对 `struct` 设置的规则进行验证(包含嵌套验证), 将最终的所有错误都返回**
* 支持对 **单个变量** 的验证, 变量可以为切片/数组/单个[int系列, float系列, bool系列, string系列]进行验证

##### 4.2 验证

###### 4.2.1 支持的验证如下:

| 标识 | 自定义 msg | 说明 |
| -------- | ---------|---------|
| required | yes |必填标识, 支持嵌套验证(结构体和结构体切片) |
| exist | yes |子对象有值才验证, 用于嵌套验证(结构体和结构体切片) |
| either | no |多选一, 即多个中必须有一个必填, 格式为 "either=xxx"(通过数据进行标识) |
| botheq | no |多都相等, 即多个中必须都相等, 格式为 "botheq=xxx"(通过数据进行标识) |
| to | yes |闭区间验证, 采用 `左右闭区间` , 格式为 "to=xxx\~xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度), 如: "to=1\~10" |
| ge | yes |大于或等于验证, 格式为 "ge=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度) |
| le | yes |小于或等于验证, 格式为: "le=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度) |
| oto | yes |开区间验证, 采用 `左右开区间` , 格式为 "oto=xxx\~xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度), 如: "oto=1\~10" |
| gt | yes |大于验证, 格式为 "gt=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度) |
| lt | yes |小于验证, 格式为: "lt=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度) |
| eq | yes |等于验证, 格式为: "eq=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度) |
| noeq | yes |不等于验证, 格式为: "noeq=xxx"(字段类型: 字符串为长度, 数字为大小, 切片为长度) |
| in | yes |指定输入选项, 格式为 "in=(xxx/xxx/xxx)", 如: "in=(1/abc/3)" |
| include | yes |指定输入包含选项, 格式为 "include=(xxx/xxx/xxx)", 如: "in=(hello/2/3)" |
| phone | yes |手机号验证 |
| email | yes |邮箱验证 |
| ipv4 | yes| ipv4 验证|
| idcard | yes |身份证号码验证 |
| year | yes |年验证 |
| year2month| yes |年月验证, 支持分割符, 默认按照"-". 验证:xxxx/xx, 格式: "year2month=/" |
| date | yes |日期验证, 支持分割符, 默认按照"-". 验证:xxxx/xx/xx, 格式: "date=/" |
| datetime | yes |时间验证, 支持分割符, 默认按照"-". 验证:xxxx/xx/xx xx\:xx\:xx, 格式: "datetime=/" |
| int | yes |整数型验证 |
| ints | yes |验证是否为多个数字. 如果输入为 string, 默认按逗号拼接进行验证; 如果为 slice/array, 会将每个值进行匹配判断 |
| float | yes |浮动数型验证 |
| re | yes |正则验证, 格式为: "re='xxx'", 如: "re='[a-z]+'" |
| unique | yes |唯一验证, 说明: 1.对以逗号隔开的字符串进行唯一验证; 2. 对切片/数组元素[int 系列, float系列, bool系列, string系列]进行唯一验证 |

* 自定义 msg 写法如下:
  + 1. 如: `required|必填`, key 为 `required`, value 为 ``, cusMsg 为 `必填`; 
  + 2. 如: `to=1~2|大于等于 1 且小于等于 2`, key 为 `to`, value 为 `1~2`, cusMsg 为 `大于等于 1 且小于等于 2`

###### 4.2.2 设置验证

*  1. 通过设置 `tag` 进行设置验证规则, 默认目标为 `valid`
*  2. 支持通过创建 `RM` 对象进行自定义设置验证规则, 其验证优先级高于 `xxx.pb.go` 里的规则,  `RM` 暂不支持嵌套

###### 4.2.3 其他

*  1. 默认按照 `tag` 进行处理, 如果设置 `RM` 对象会以此规则为准
*  2. 如果验证方法没有实现的, 可以调用 `SetCustomerValidFn` 自定义
*  3. 使用的可以参考 `example_test.go` 和 `valid_test.go`

#### 5 使用示例:

*  `proto` 内容如下:  

```
message Man {
    string name = 1; // 姓名 @tag valid:"required,to=1~3"
    int32 age = 2; // 年龄 @tag valid:"to=1~150"
}
```

*  **注:** 编写 `xxx.proto` 时, 需要加将 `@tag xxx` 放到注释的最后面
* 执行命令: `inject_tool.sh xxx.proto` 生成 `pd` 内容如下:  

```
type Man struct {
    state protoimpl.MessageState
    sizeCache protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" valid:"required,to=1~3"` // 姓名 @tag valid:"required,to=1~3"
    Age int32 `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty" valid:"to=1~150"` // 年龄 @tag valid:"to=1~150"
}
```

* 代码里的使用  

```
m := &test.Man{
    Name: "xue",
    Age: -1,
}
fmt.Println(ValidateStruct(m))

// Output: "Man.Age" input "-1" is size less than 1
```

#### 最后

* 欢迎大佬们指正, 同时也希望大佬给 **star**，[to gitee](https://gitee.com/xuesongtao/protoc-go-valid)
