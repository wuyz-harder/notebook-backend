package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Device struct {
	Id     int    `json:"id" `
	Name   string `json:"name"`
	Owner  int    `json:"owner" `
	Code   string `json:"code"`
	Online int    `json:"online"`
	Size   string `json:"size"`
	// 后面这个是设置自动生成时间，默认是当前值
	CreateTime time.Time `gorm:"column:create_time;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create"`
}

// 根据用户id查询设备，待会要联合查询join
func (dev *Device) GetDevices() ([]Device, error) {
	devices := []Device{}
	err := Db.Find(&devices).Error
	return devices, err
}

// 根据用户id查询设备，待会要联合查询join
func (dev *Device) GetDeviceByUserId(id int) (interface{}, error) {
	var user User
	// 查询userid下面的设备
	err := Db.Preload("Devices").Find(&user, id).Error
	// err := Db.Where(&Device{Owner: id}).Find(&devices).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return user, err
}

// 根据用户id查询设备，待会要联合查询join
func (dev *Device) GetDeviceById(id int) (interface{}, error) {
	device := Device{}
	err := Db.Where(&Device{Id: id}).Find(&device).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return device, err
	}
	return device, nil
}

func (dev *Device) DeleteDeviceById(id int) (Device, error) {

	device := Device{}

	err := Db.Delete(&Device{Id: id}).Error

	fmt.Print(err)
	if err != nil {
		Db.Rollback()
		return device, err
	}
	return device, nil
}

func (dev *Device) AddDevice(device *Device) (Device, error) {

	err := Db.Transaction(func(tx *gorm.DB) error {
		creErr := tx.Create(device).Error
		if creErr != nil {
			return creErr
		}
		return nil
	})

	if err != nil {
		return *device, err
	}
	return *device, nil
}

// 更新设备信息
// 由于map只会更新非空的字段，如果有非空的过来就不更新字段值，因此改用map
// 通过结构体变量更新字段值, gorm库会忽略零值字段。就是字段值等于0, nil, "", false这些值会被忽略掉，不会更新。如果想更新零值，可以使用map类型替代结构体。
func (dev *Device) UpdateDevice(device map[string]interface{}) (interface{}, error) {

	id := device["id"]
	resDevice := Device{}
	err := Db.Transaction(func(tx *gorm.DB) error {
		// 只会更新有值的字段

		updateErr := tx.Model(&Device{}).Where("id=?", id).Updates(device).Error
		return updateErr
	})

	if err != nil {

		return device, err
	}
	// 类型转换
	Db.Model(&Device{}).Where("id=?", id).Find(&resDevice)
	return resDevice, nil
}
