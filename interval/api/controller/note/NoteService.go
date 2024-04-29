package note

import (
	"encoding/json"
	"strconv"

	"github.com/wuyz-harder/notebook-backend/common"
	"github.com/wuyz-harder/notebook-backend/interval/api/models"
	"github.com/wuyz-harder/notebook-backend/utils"

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

/*
删除笔记
*/
func DeleteNote(context *gin.Context) {
	var note models.Note
	// 获取path参数
	if context.Param("id") == "" {
		context.Error(common.NewError(400, 200400, "参数信息错误"))
		return
	}
	id, _ := strconv.Atoi(context.Param("id"))
	note.ID = id
	// 解析用户信息
	token := context.GetHeader("token")
	_, claim, tokenErr := utils.ParseToken(token)
	if tokenErr != nil {
		context.Error(common.NewError(401, 200401, "用户信息错误"))
		return
	}
	ownerID, ownerErr := note.GetNoteOwnerID()

	if ownerErr != nil {
		context.Error(common.NewError(400, 200400, ownerErr.Error()))
		return
	}
	// 如果不是自己的，那么权限也会没有
	if ownerID != claim.UserID {
		context.Error(common.NewError(400, 200400, "没有权限，需要本人"))
		return
	}
	// 执行删除
	err := note.Delete()
	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "success", map[string]interface{}{"id": id}))
		return
	}
}

// 分享笔记

func ShareNote(context *gin.Context) {
	var note models.Note

	body, _ := context.GetRawData()
	json.Unmarshal(body, &note)
	// 获取path参数
	if note.ID == 0 {
		context.Error(common.NewError(400, 200400, "参数信息错误,id错误"))
		return
	}
	// 解析用户信息
	token := context.GetHeader("token")
	_, claim, tokenErr := utils.ParseToken(token)
	if tokenErr != nil {
		context.Error(common.NewError(401, 200401, "用户信息错误"))
		return
	}
	ownerID, ownerErr := note.GetNoteOwnerID()
	if ownerErr != nil {
		context.Error(common.NewError(400, 200400, ownerErr.Error()))
		return
	}
	// 如果不是自己的，那么权限也会没有
	if ownerID != claim.UserID {
		context.Error(common.NewError(400, 200400, "没有权限，需要本人"))
		return
	}
	// 获取分享链接
	url, err := note.Share()
	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "success", map[string]interface{}{"id": note.ID, "url": url}))
		return
	}
}

func GetShareNote(context *gin.Context) {
	shareUrl := context.Query("shareUrl")
	var note models.Note
	var author models.User
	// 获取path参数
	if shareUrl == "" {
		context.Error(common.NewError(400, 200400, "shareUrl参数错误"))
		return
	}
	res, err := note.GetNotebyShareUrl(shareUrl)
	author.ID = res.User
	user, _ := author.GetUserByID()

	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "success", map[string]interface{}{"note": res, "author": user}))
		return
	}
}
