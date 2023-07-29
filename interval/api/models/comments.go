package models

import (
	"github.com/jinzhu/gorm"
)

// 外键在gorm里面要定义两个字段：
// 1、跟数据库对应的外键id字段，只保存数值，比如下面的Commenter
// 2、与外键实体对应的类型字段，实体字段后面的定义要加上gorm:"foreignkey:Commenter"指定
// 注意，preload的时候写的是实体字段的名称，比如下面的CommenterInfo

type Comment struct {
	ID int `json:"id" `
	// 评论内容
	Content string `json:"content"`
	// 笔记外键
	Note          int   `json:"note"`
	Commenter     int   `json:"commenterID"`
	CommenterInfo *User `gorm:"foreignkey:Commenter" json:"commenterInfo"` // 关联外键字段，表示该用户的父用户
	// *int代表外键可以为空，自引用的时候一定要加*号
	ReplyID      *int     `json:"replyID"`
	ReplyComment *Comment `gorm:"foreignkey:ReplyID"` // 关联外键字段，表示该用户的父用户

	// 后面这个是设置自动生成时间，默认是当前值
	CreateTime *LocalTime `gorm:"column:create_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
}

// 创建评论
func (comment *Comment) Create() (int, error) {
	var id int = 0
	err := Db.Transaction(func(tx *gorm.DB) error {
		cerr := tx.Create(&comment).Error
		id = comment.ID
		return cerr

	})
	return id, err
}

// 获取某个笔记下的所有评论

func (comment *Comment) GetALL(noteID int) ([]Comment, error) {
	comments := []Comment{}
	err := Db.Preload("CommenterInfo").Preload("ReplyComment").Where(&Comment{Note: noteID}).Find(&comments).Error
	return comments, err
}
