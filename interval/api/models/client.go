package models

import (
	"GetHotWord/utils"
	"errors"

	"github.com/jinzhu/gorm"
)

type WsClient struct {
	ID int `json:"id"`

	User        int    `json:"user"`
	UserInfo    *User  `gorm:"foreignkey:User" json:"userInfo"` // 关联外键字段，表示该用户的父用户
	WebsocketID string `json:"websocketID"`
	Online      int    `json:"online"`

	// 后面这个是设置自动生成时间，默认是当前值
	LastActive *LocalTime `gorm:"column:last_active;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
}

func (client *WsClient) GetClientByUser() error {

	err := Db.Where(&WsClient{User: client.User}).Preload("UserInfo").Find(&client).Error
	if err != nil {
		// 如果没有的话就创建
		if errors.Is(err, gorm.ErrRecordNotFound) {
			client.Create()
			return nil
		}
		return err
	} else {
		return nil
	}
}

// 获取在线的设备或者用户
func (client *WsClient) GetClientOnline() ([]WsClient, error) {
	var res []WsClient
	err := Db.Preload("UserInfo").Where(&WsClient{Online: 1}).Find(&res).Error
	if err != nil {
		return res, err
	} else {
		return res, nil
	}
}

// 获取在线的设备或者用户
func (client *WsClient) Create() error {
	// 生成随机码
	client.WebsocketID = utils.GenerateClientCode()
	// 在线，用户信息在调用这个函数前就已经设置了
	client.Online = 1
	err := Db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(client).Error
	})
	return err
}

// 保存设备
func (client *WsClient) Save() error {

	err := Db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(client).Error
	})
	return err
}
