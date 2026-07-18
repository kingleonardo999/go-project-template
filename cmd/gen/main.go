package main

import (
	"log"

	"go-project-template/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.Load("config/app.yaml"); err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	// 创建生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:      "internal/repository/db/dao", // DAO 代码输出路径
		ModelPkgPath: "internal/model",             // model 包路径
		Mode:         gen.WithDefaultQuery | gen.WithoutContext,

		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(db)

	// 生成所有表
	g.ApplyBasic(g.GenerateAllTable()...)

	// 自定义查询方法示例（按需启用）：
	// 1. 在 internal/model/querier.go 中定义接口
	// 2. 用 g.ApplyInterface 绑定到对应表
	//
	// 例：g.ApplyInterface(func(model.UserQuerier) {}, g.GenerateModel("users"))
	// 例：g.ApplyInterface(func(model.OrderQuerier) {}, g.GenerateModel("orders"))

	g.Execute()
}
