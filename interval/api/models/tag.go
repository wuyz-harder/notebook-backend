package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	ID    uint   `json:"id" `
	Name  string `json:"name" `
	Count int    `json:"count" `
	Users []User `gorm:"many2many:user_tag_tables"`
}

func (tag *Tag) Exist(fingTag *Tag) (bool, Tag) {
	var resTag Tag
	err := Db.Model(&Tag{}).Where("name=?", fingTag.Name).Find(&resTag).Error
	if err == nil {
		return true, resTag
	} else {
		// 如果不保存就返回为空
		return false, resTag
	}
}

// 增加
func (tag *Tag) Insert(name string) error {
	err := Db.Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&Tag{}).Create(&Tag{Name: name}).Error

		return err
	})
	return err
}

func (tag *Tag) Delete(name string) error {
	err := Db.Transaction(func(tx *gorm.DB) error {

		err := tx.Model(&Tag{}).Delete(&Tag{Name: name}).Error

		return err
	})
	return err

}

func (tag *Tag) GetAll() ([]Tag, error) {
	var res []Tag
	err := Db.Model(&Tag{}).Preload("Users").Find(&res).Error
	return res, err

}

// 标签数量+1
func (tag *Tag) AddCount() error {
	// var  rwLock sync.RWMutex
	// rwLock.
	err := Db.Transaction(func(tx *gorm.DB) error {
		// 加锁操作
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Find(tag).Error; err != nil {
			tx.Rollback()
			return err
		}
		tag.Count = tag.Count + 1

		upErr := tx.Save(tag).Error
		return upErr
	})

	return err
}
