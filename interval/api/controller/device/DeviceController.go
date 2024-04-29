package device

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wuyz-harder/notebook-backend/common"
	"github.com/wuyz-harder/notebook-backend/interval/api/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var MyDevice models.Device

// 查询属于某个用户的设备
func GetDeviceByUserId(context *gin.Context) {

	id := context.Query("id")
	if id == "" {
		context.Error(common.NewError(400, 200400, "id参数缺失"))
		return
	}
	nid, err := strconv.Atoi(id)
	if err != nil {
		context.Error(common.NewError(400, 200400, "id参数不是整数"))
		return
	}
	devices, queryErr := MyDevice.GetDeviceByUserId(nid)
	if queryErr != nil {
		context.Error(common.NewError(400, 200400, queryErr.Error()))
		return
	}

	context.JSON(200, gin.H{
		"res": devices,
	})
}

// 查询所有设备
func GetDevices(context *gin.Context) {

	devices, err := MyDevice.GetDevices()
	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	}
	context.JSON(200, gin.H{
		"res": devices,
	})
}

// 根据id获取设备
func GetDeviceById(context *gin.Context) {
	id := context.Query("id")
	if id == "" {
		zap.L().Info("缺少id参数")
		context.JSON(
			400, gin.H{
				"res":    "fail",
				"reason": "缺少id参数",
			})

		return
	}
	nid, err := strconv.Atoi(id)
	if err != nil {
		context.Error(common.NewError(400, 200400, "id参数不是整数"))
		zap.L().Info("id参数不是整数")
		return
	}

	device, finderr := MyDevice.GetDeviceById(nid)
	if finderr != nil {
		context.JSON(200,
			gin.H{
				"data": "",
			})
		return
	}
	context.JSON(200,
		gin.H{
			"data": device,
		},
	)
}

// 删除属于某个用户的设备
func DeleteDeviceById(context *gin.Context) {

	id := context.Query("id")
	if id == "" {
		zap.L().Info("id参数缺失")
		context.Error(common.NewError(400, 200400, "id参数缺失"))
		return
	}
	nid, err := strconv.Atoi(id)
	if err != nil {
		zap.L().Info("id参数不是整数")
		context.Error(common.NewError(400, 200400, "id参数不是整数"))
		return
	}
	_, delErr := MyDevice.DeleteDeviceById(nid)
	if delErr != nil {
		context.Error(common.NewError(400, 200400, delErr.Error()))
		return
	}
	context.JSON(200, gin.H{
		"res": "success",
	})
}

func AddDevice(context *gin.Context) {

	var device models.Device
	data, err := context.GetRawData()
	if err != nil {
		context.Error(common.NewError(500, 200400, "后台出错了"))
		return
	}
	// json反序列化，以一个map来进行接收
	json.Unmarshal(data, &device)
	fmt.Println(data)
	if device.Owner == 0 {
		zap.L().Error("owner参数错误")
		context.Error(common.NewError(400, 200400, "owner参数错误"))
		return
	}
	if device.Name == "" {
		zap.L().Error("name参数错误")
		context.Error(common.NewError(400, 200400, "name参数错误"))
		return
	}

	if device.Code == "" {
		zap.L().Error("code参数错误")
		context.Error(common.NewError(400, 200400, "code参数错误"))
		return
	}
	_, addErr := MyDevice.AddDevice(&device)
	if addErr != nil {
		context.Error(common.NewError(400, 200400, addErr.Error()))
	} else {
		context.JSON(200, gin.H{"mes": "success"})
	}

}

// 修改设备信息
func UpdateDevice(context *gin.Context) {

	var device map[string]interface{}
	data, err := context.GetRawData()
	if err != nil {
		context.Error(common.NewError(500, 200400, "后台出错了"))
		return
	}
	// json反序列化，以一个map来进行接收
	json.Unmarshal(data, &device)
	if _, IDok := device["id"]; !IDok {
		zap.L().Error("id参数错误")
		context.Error(common.NewError(400, 200400, "id参数错误"))
		return
	}
	updatedDev, addErr := MyDevice.UpdateDevice(device)
	if addErr != nil {
		context.Error(common.NewError(400, 200400, addErr.Error()))
	} else {
		context.JSON(200, gin.H{"mes": "success", "data": updatedDev})
	}

}
