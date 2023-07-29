package models

import (
	"fmt"

	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/exp/rand"
)

// gorm时间格式化
type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

type Note struct {
	ID       int    `json:"id" `
	Title    string `json:"title"`
	Size     string `json:"size"`
	ShareUrl string `json:"shareUrl"`
	// 文章内容必须是长文本
	Content string `json:"content" gorm:"type:LONGTEXT"`
	// 后面这个是设置自动生成时间，默认是当前值
	CreateTime *LocalTime `gorm:"column:create_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
	// 后面这个是设置自动生成时间，默认是当前值，在update的时候自动更新
	UpdateTime *LocalTime `gorm:"column:update_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:update"`
	// omitempty前端可以传过来但是后端隐射回去的时候不用显示
	User int `json:"user,omitempty"`
}

func (note *Note) Create() (int, error) {
	var id int = 0
	// 创建或者更新
	err := Db.Transaction(func(tx *gorm.DB) error {
		if note.ID == 0 {
			cERR := tx.Create(&note).Error
			id = note.ID
			return cERR
		} else {
			saveERR := tx.Save(&note).Error
			id = note.ID
			return saveERR
		}

	})

	return id, err
}

// 获取自己的所有笔记以及数量
func (note *Note) GetALL(userID int) ([]Note, int, error) {
	res := []Note{}
	num := 0
	checkRes := Db.Where(&Note{User: userID}).Find(&res)
	//  查询数量
	checkRes.Count(&num)
	return res, num, checkRes.Error
}

func (note *Note) GetNoteOwnerID() (int, error) {
	res := Note{}
	err := Db.Where(&Note{ID: note.ID}).Find(&res).Error
	if err != nil {
		return 0, err
	}
	return res.User, nil
}

func (note *Note) Delete() error {
	err := Db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(note).Error
	})
	return err
}

// ShareUrl为空的话则生成，有的话则直接获取
func (note *Note) Share() (string, error) {

	fErr := Db.Where(&Note{ID: note.ID}).Find(&note).Error
	if fErr != nil {
		return "", fErr
	}
	// 如果本来还没有的话就要生成
	if note.ShareUrl == "" {
		// 生成十个随机字符串
		var str = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
		usr := ""
		for i := 1; i < 10; i++ {
			usr = usr + string(str[rand.Intn(35)])
		}
		note.ShareUrl = usr
		// 先保存
		saveErr := Db.Transaction(func(tx *gorm.DB) error {
			return tx.Save(&note).Error
		})
		if saveErr != nil {
			return "", saveErr
		}
		return usr, nil
	} else {
		// 如果之前有的话直接返回
		return note.ShareUrl, nil
	}

}

// ShareUrl为空的话则生成，有的话则直接获取
func (note *Note) GetNotebyShareUrl(shareUrl string) (Note, error) {
	res := Note{}
	fErr := Db.Where(&Note{ShareUrl: shareUrl}).Find(&res).Error
	if fErr != nil {
		return Note{}, fErr
	} else {
		return res, nil
	}

}
