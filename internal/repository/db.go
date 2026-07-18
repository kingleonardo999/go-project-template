package repository

import "gorm.io/gorm"

// DB 全局数据库实例，service 层通过它访问数据
var DB *gorm.DB

// SetDB 初始化全局 DB 实例
func SetDB(db *gorm.DB) {
	DB = db
}

// GetDB 获取 DB 实例，方便测试时替换
func GetDB() *gorm.DB {
	return DB
}
