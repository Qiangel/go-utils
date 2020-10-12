package log

import (
	"context"
	"fmt"
	"testing"

	"github.com/Qiangel/go-utils/mysql"
)

func TestPrintLog(t *testing.T) {
	ctx := context.TODO()
	ctx = NewTraceIDContext(ctx, "1213dsd")
	Init() /**/
	Info(ctx, "xxx1")
	Info(ctx, "xxx2")
}

func TestHooker_Write(t *testing.T) {
	cnf := mysql.Config{
		Username: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     "3306",
		Dbname:   "test", //日志存储的数据库信息
	}
	err := InitHook(cnf)
	fmt.Println(err)
	Init(
		WithEnableSaveFile(false), //关闭保存文件
		WithEnableSaveDB(true),    //打开保存数据库
		WithFilePath("./logs/"),   //WithEnableSaveFile的参数相关
		WithVersion("v1.1.1"),     // 版本号
		WithMaxAge(1),             //WithEnableSaveFile的参数相关
		WithMaxSize(1),            //WithEnableSaveFile的参数相关
		WithMaxBackups(2),         //WithEnableSaveFile的参数相关
	)
	ctx := context.Background()
	ctx = NewTraceIDContext(ctx, "12125484512")
	ctx = NewUserIDContext(ctx, "nameduserid")
	for i := 0; i <= 0; i++ {
		Info(ctx, "this is some log", i)
	}
}
