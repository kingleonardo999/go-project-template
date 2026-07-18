// Package model 存放 GORM Gen 生成的数据库实体和自定义查询接口。
//
// 自定义查询接口定义在此处，gen 会自动生成实现。
// 方法上的注释就是 gen 执行的 SQL，@@table 会被替换为实际表名。
package model

import "context"

// UserQuerier 示例：用户表的自定义查询方法
type UserQuerier interface {
	// SELECT * FROM @@table WHERE email = @email LIMIT 1
	GetByEmail(ctx context.Context, email string) (map[string]any, error)

	// SELECT * FROM @@table WHERE status = @status ORDER BY created_at DESC
	GetByStatus(ctx context.Context, status int) ([]map[string]any, error)

	// 分页查询
	// SELECT * FROM @@table ORDER BY id LIMIT @size OFFSET (@page - 1) * @size
	FindByPage(ctx context.Context, page, size int) ([]map[string]any, int64, error)
}

// OrderQuerier 示例：订单表的自定义查询方法
type OrderQuerier interface {
	// SELECT * FROM @@table WHERE user_id = @userID ORDER BY created_at DESC
	GetByUserID(ctx context.Context, userID int64) ([]map[string]any, error)

	// SELECT * FROM @@table WHERE status = @status AND created_at BETWEEN @start AND @end
	FindByStatusAndCreatedAtBetween(ctx context.Context, status int, start, end string) ([]map[string]any, error)
}
