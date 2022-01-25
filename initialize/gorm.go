package initialize

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMySQL() *gorm.DB {
	// 打开MySQL日志文件
	path, err := os.Executable()
	if err != nil {
		panic("初始化失败：可执行程序路径获取失败")
	}
	path = filepath.Dir(path)
	path = filepath.Join(path, "phoenix-mysql.log")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		panic("初始化失败：打开MySQL日志文件失败")
	}
	// 定制MySQL的Logger
	newLogger := logger.New(
		log.New(file, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)
	// 读取配置数据
	addr := global.VP.GetString("database.ip")
	port := global.VP.GetString("database.port")
	user := global.VP.GetString("database.user")
	password := global.VP.GetString("database.password")
	database := global.VP.GetString("database.database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, addr, port, database)
	// 连接MySQL数据库
	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("初始化失败：连接MySQL数据库失败")
	}
	// 更新MySQL数据库内容
	db.Set("gorm:table_options", "ENGINE=InnoDB")
	db.AutoMigrate(
		&model.User{},
		&model.Problem{},
		&model.Tutorial{},
		&model.Competition{},
		&model.Organization{},
		&model.Comment{},
		&model.Post{},
	)
	return db
}
