package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// 定义双方信息的的载体
type Message struct {
	ID int `json:"id" `
	// 评论内容
	Content string `json:"content"`
	// 笔记外键
	From int `json:"from"`
	To   int `json:"to"`
	// *int代表外键可以为空，自引用的时候一定要加*号
	FromUserInfo *User `gorm:"foreignkey:From" json:"fromUserInfo"` // 关联外键字段，表示该用户的父用户
	// *int代表外键可以为空，自引用的时候一定要加*号
	ToUserInfo *User `gorm:"foreignkey:To" json:"toUserInfo"` // 关联外键字段，表示该用户的父用户
	// 后面这个是设置自动生成时间，默认是当前值
	CreateTime *LocalTime `gorm:"column:create_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
}

// 创建消息
func (message *Message) Create() (Message, error) {
	var preMessage Message

	err := Db.Transaction(func(tx *gorm.DB) error {
		cerr := tx.Create(&message).Error
		return cerr

	})
	if err == nil {
		Db.Preload("FromUserInfo").Preload("ToUserInfo").First(&preMessage, message.ID)
		return preMessage, nil
	}
	return preMessage, err
}

// 获取某两人的所有消息

func (message *Message) GetALL(toUserID int) ([]Message, error) {
	// 需要查询双方的
	fmt.Println(message)
	messages := []Message{}
	nums := []int{message.From, message.To}
	fmt.Println(nums)

	err := Db.
		Preload("FromUserInfo").
		Preload("ToUserInfo").
		Where(message).
		Or(&Message{From: message.To, To: message.From}).
		Find(&messages).
		Error
	return messages, err
}
