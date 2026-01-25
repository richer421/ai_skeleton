package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:@tcp(127.0.0.1:3306)/ai_skeleton?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 创建生成器
	g := gen.NewGenerator(gen.Config{
		OutPath: "./repository/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	// 生成所有表的代码
	g.ApplyBasic(
		// 在这里添加需要生成的模型
		// g.GenerateModel("users"),
	)

	// 执行生成
	g.Execute()

	log.Println("Code generation completed!")
}
