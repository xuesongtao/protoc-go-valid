# proto 中注入 tag
#### 1. 介绍
- 1. 对 `xxx.proto` 文件注入 tag 
- 2. 通过验证器对内容进行验证
- **注:** 必须先按照 `protoc`


#### 2. 注入工具使用
- 1. 先下载本项目: `go get -u gitee.com/xuesongtao/protoc-go-valid`
- 2. `protoc-go-valid` 支持 3 种操作, 如下: 
	- 2.1 `protoc-go-valid -d="待注入的目录"`
	- 2.2 `protoc-go-valid -p="匹配模式"`
	- 2.3 `protoc-go-valid -f="单个待注入的文件"`
	

#### 3. 工具补充
- 1. protoc-go-valid -h` 可以通过这个查询示例
- 2. 由于此操作是先执行 `protoc` 才能注入的, 当前项目有个 shell 脚本, 可以把项目中的`protoc.sh`加到环境变量, **说明:** 如果为 **windows** 需要使用 `powershell` 来执行
- 3. 修改 `protoc.sh` 中 `proto` 和 `pb` 的目录, 相当于应用的目录; 如本项目, 修改如下下:
```
outPdProjectPath="test" # pb 放入的项目路径
protoFileDirName="test" # proto 存放的目录
```



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