### 1开启切割和存储到文件
```go
//初始化
Init(
		WithEnableSaveFile(true), //关闭保存文件
		WithFilePath("./logs/"),   //WithEnableSaveFile的参数相关
		WithMaxAge(1),             //WithEnableSaveFile的参数相关
		WithMaxSize(1),            //WithEnableSaveFile的参数相关
		WithMaxBackups(2),         //WithEnableSaveFile的参数相关
	)
//使用
ctx := context.Background()
ctx = NewTraceIDContext(ctx, "12125484512")
ctx = NewUserIDContext(ctx, "nameduserid")
Info(ctx, "this is some log")
```


### 2开启存储到数据库 可自动创建数据库模型
```go
//初始化
cnf := mysql.Config{
		Username: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     "3306",
		Dbname:   "test", //日志存储的数据库信息
	}
err := InitHook(cnf)
Init(WithEnableSaveDB(true))//打开保存数据库
//使用
ctx := context.Background()
ctx = NewTraceIDContext(ctx, "12125484512")
ctx = NewUserIDContext(ctx, "nameduserid")
Info(ctx, "this is some log")
```