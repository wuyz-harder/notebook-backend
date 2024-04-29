package star

import (
	"encoding/json"
	"fmt"

	"github.com/wuyz-harder/notebook-backend/common"
	"github.com/wuyz-harder/notebook-backend/interval/api/models"

	"github.com/gin-gonic/gin"
)

func CreateStar(context *gin.Context) {

	var star models.Star
	body, err := context.GetRawData()
	userID, _ := context.Get("userID")
	if err != nil {
		context.Error(common.NewError(400, 200400, "参数错误"))
		return
	}
	jsonErr := json.Unmarshal(body, &star)
	star.User = userID.(int)
	fmt.Println(star)
	if jsonErr != nil {
		context.Error(common.NewError(400, 200400, "参数错误"))
		return
	}
	if star.Note == 0 {
		context.Error(common.NewError(400, 200400, "note参数为空"))
		return
	}
	id, cErr := star.Create()
	if cErr != nil {
		context.Error(common.NewError(400, 200400, cErr.Error()))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "创建成功", map[string]interface{}{"id": id}))
		return
	}

}

func GetStars(context *gin.Context) {

	var star models.Star
	userID, _ := context.Get("userID")
	star.User = userID.(int)
	res, cErr := star.GetALL()
	if cErr != nil {
		context.Error(common.NewError(400, 200400, cErr.Error()))
		return
	} else {
		context.JSON(200, common.NewResp(200, 200200, "sueecss", map[string]interface{}{"stars": res}))
		return
	}

}
