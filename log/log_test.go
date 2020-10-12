package log

import (
	"context"
	"fmt"
	"testing"

	"github.com/Qiangel/go-utils/mysql"
	"github.com/gin-gonic/gin"
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

func TestGin(t *testing.T) {
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
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		/*中间件位置s  */
		ctx := NewUserIDContext(c.Request.Context(), "this is mock userID")
		ctx = NewTraceIDContext(ctx, "this is random mock traceID") //随机traceid
		c.Request = c.Request.WithContext(ctx)
		/*中间件位置e*/

		//Info(c, "this is some log1")
		Info(ctx, "this is some log2")

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
