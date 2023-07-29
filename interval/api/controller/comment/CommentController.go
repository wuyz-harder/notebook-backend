package comment

import (
	"GetHotWord/common"
	"GetHotWord/interval/api/models"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建note接口
func CreateComment(context *gin.Context) {
	var comment models.Comment
	// 通过请求头获取id,设置评论人
	userID, _ := context.Get("userID")
	comment.Commenter = userID.(int)

	body, getDataErr := context.GetRawData()
	if getDataErr != nil {
		context.Error(common.NewError(400, 200400, "参数获取错误"))
		return
	}
	parseErr := json.Unmarshal(body, &comment)
	if parseErr != nil {
		context.Error(common.NewError(400, 200400, "参数错误"))
		return
	}
	// 计算字节大小
	id, createErr := comment.Create()
	if createErr != nil {
		context.Error(common.NewError(400, 200400, "创建错误"))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "创建成功", map[string]interface{}{"id": id}))
		return
	}

}

// 创建note接口
func GetComments(context *gin.Context) {
	var comment models.Comment
	noteID, err := strconv.Atoi(context.Query("noteID"))
	if err != nil {
		context.Error(common.NewError(400, 200400, "参数获取错误"))
		return
	}
	res, getErr := comment.GetALL(noteID)
	if getErr != nil {
		context.Error(common.NewError(500, 200500, "查询错误"))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "success", map[string]interface{}{"comments": res}))
		return
	}

}
