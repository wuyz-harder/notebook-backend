package tag

import (
	"encoding/json"

	"github.com/wuyz-harder/notebook-backend/common"
	"github.com/wuyz-harder/notebook-backend/interval/api/models"

	"github.com/gin-gonic/gin"
)

var MyTag models.Tag

type TagArg struct {
	Name string
}

// 获取所有标签
func GetAll(ctx *gin.Context) {

	res, err := MyTag.GetAll()
	if err != nil {
		ctx.Error(common.NewError(500, 200500, err.Error()))
		return
	}
	ctx.JSON(200, gin.H{
		"msg":  "success",
		"data": res,
	})
}

// 增加标签
func Add(ctx *gin.Context) {
	var tagArg TagArg
	body, err := ctx.GetRawData()
	if err != nil {
		ctx.Error(common.NewError(500, 200500, err.Error()))
		return
	}
	parseErr := json.Unmarshal(body, &tagArg)
	if parseErr != nil {
		ctx.Error(common.NewError(500, 200500, parseErr.Error()))
		return
	}
	if tagArg.Name == "" {
		ctx.Error(common.NewError(400, 200400, "lost parameter name"))
		return
	}
	insertErr := MyTag.Insert(tagArg.Name)
	if insertErr == nil {
		ctx.JSON(200, gin.H{
			"msg":  "success",
			"data": tagArg.Name,
		})

	} else {
		ctx.Error(common.NewError(500, 200500, insertErr.Error()))
		return
	}
}

// 删除标签
func Delete(ctx *gin.Context) {
	var tagArg TagArg
	body, err := ctx.GetRawData()
	if err != nil {
		ctx.Error(common.NewError(500, 200500, err.Error()))
		return
	}
	parseErr := json.Unmarshal(body, &tagArg)
	if parseErr != nil {
		ctx.Error(common.NewError(500, 200500, parseErr.Error()))
		return
	}
	if tagArg.Name == "" {
		ctx.Error(common.NewError(400, 200400, "lost parameter name"))
		return
	}
	insertErr := MyTag.Delete(tagArg.Name)
	if insertErr == nil {
		ctx.JSON(200, gin.H{
			"msg":  "delete success",
			"data": tagArg.Name,
		})

	} else {
		ctx.Error(common.NewError(500, 200500, insertErr.Error()))
		return
	}

}
