package message

import (
	"GetHotWord/common"
	"GetHotWord/interval/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建note接口
func CreateMessage(fromUserID int, toUserID int, content string) (models.Message, error) {
	var message models.Message
	// 计算字节大小
	res, createErr := message.Create()
	if createErr != nil {
		return res, createErr
	} else {
		return res, nil
	}
}

// 创建note接口
func GetMessages(context *gin.Context) {
	// 获取from userID

	fromUserID, _ := context.Get("userID")
	toUserID, err := strconv.Atoi(context.Query("toUserID"))
	message := models.Message{
		From: fromUserID.(int),
		To:   toUserID,
	}
	if err != nil {
		context.Error(common.NewError(400, 200400, "参数获取错误"))
		return
	}
	res, getErr := message.GetALL(toUserID)
	if getErr != nil {
		context.Error(common.NewError(500, 200500, "查询错误"))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "success", map[string]interface{}{"messages": res}))
		return
	}

}
