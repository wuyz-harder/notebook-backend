package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// gorm时间格式化
type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

type Note struct {
	ID    int    `json:"id" `
	Title string `json:"title"`
	Size  string `json:"size"`
	// 文章内容必须是长文本
	Content string `json:"content" gorm:"type:LONGTEXT"`
	// 后面这个是设置自动生成时间，默认是当前值
	CreateTime LocalTime `gorm:"column:create_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	// omitempty前端可以传过来但是后端隐射回去的时候不用显示
	User int `json:"user,omitempty"`
}

func (note *Note) Create() (int, error) {
	var id int = 0
	err := Db.Transaction(func(tx *gorm.DB) error {

		cERR := tx.Create(&note).Error
		fmt.Println(note.ID)
		id = note.ID
		return cERR
	})

	return id, err
}

// 获取所有笔记以及数量
func (note *Note) GetALL(userID int) ([]Note, int, error) {
	res := []Note{}
	num := 0
	fmt.Println(userID)
	checkRes := Db.Where(&Note{User: userID}).Find(&res)
	//  查询数量
	checkRes.Count(&num)
	return res, num, checkRes.Error
}
