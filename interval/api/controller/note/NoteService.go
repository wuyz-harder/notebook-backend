package note

import (
	"GetHotWord/common"
	"GetHotWord/interval/api/models"
	"GetHotWord/utils"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建note接口
func CreateNote(context *gin.Context) {

	var note models.Note
	body, getDataErr := context.GetRawData()
	if getDataErr != nil {
		context.Error(common.NewError(400, 200400, "参数获取错误"))
		return
	}
	parseErr := json.Unmarshal(body, &note)
	if parseErr != nil {
		context.Error(common.NewError(400, 200400, "参数错误"))
		return
	}
	// 计算字节大小
	note.Size = strconv.Itoa(len(note.Content))
	id, createErr := note.Create()
	if createErr != nil {
		context.Error(common.NewError(400, 200400, "创建错误"))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "创建成功", map[string]interface{}{"id": id}))
		return
	}

}

func AllNote(context *gin.Context) {
	var note models.Note
	token := context.GetHeader("token")
	_, claim, tokenErr := utils.ParseToken(token)
	if tokenErr != nil {
		context.Error(common.NewError(401, 200401, "用户信息错误"))
		return
	}
	res, num, err := note.GetALL(claim.UserID)

	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "success", map[string]interface{}{"num": num, "notes": res}))
		return
	}

}
