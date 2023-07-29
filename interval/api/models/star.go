package models

import (
	"github.com/jinzhu/gorm"
)

// 外键在gorm里面要定义两个字段：
// 1、跟数据库对应的外键id字段，只保存数值，比如下面的Commenter
// 2、与外键实体对应的类型字段，实体字段后面的定义要加上gorm:"foreignkey:Commenter"指定
// 注意，preload的时候写的是实体字段的名称，比如下面的CommenterInfo

type Star struct {
	ID int `json:"id" `
	// 笔记外键
	Note     int  `json:"note"`
	NoteInfo Note `gorm:"foreignkey:Note" json:"noteInfo"` // 关联外键字段，表示该用户的父用户
	User     int  `json:"user"`
	UserInfo User `gorm:"foreignkey:User" json:"userInfo"` // 关联外键字段，表示该用户的父用户

	// 后面这个是设置自动生成时间，默认是当前值
	CreateTime LocalTime `gorm:"column:create_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
}

// 创建评论
func (star *Star) Create() (int, error) {
	var id int = 0
	err := Db.Transaction(func(tx *gorm.DB) error {
		cerr := tx.Create(&star).Error
		id = star.ID
		return cerr

	})
	return id, err
}

// 获取某个用户所有的笔记

func (star *Star) GetALL() ([]Star, error) {
	Stars := []Star{}
	// 需要预加载note的信息
	err := Db.Preload("NoteInfo").Where(&Star{User: star.User}).Find(&Stars).Error
	return Stars, err
}
