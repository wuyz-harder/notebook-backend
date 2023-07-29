package user

import (
	"GetHotWord/common"
	"GetHotWord/interval/api/models"
	"GetHotWord/utils"
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NamePW struct {
	Password string "json:password"
	Name     string "json:name"
}
type NamePWE struct {
	Password string "json:password"
	Name     string "json:name"
	Email    string "json:email"
}

var MyUser models.User

// 创建用户接口

func CreateUser(context *gin.Context) {
	// 获取post方法的表单请求参数

	var req NamePWE
	body, bderr := context.GetRawData()
	if bderr != nil {
		context.Error(common.NewError(400, 200400, "参数解析错误"))
		return
	}
	json.Unmarshal(body, &req)

	success := common.Success{
		Mes: "success",
	}

	if req.Name == "" {
		context.Error(common.NewError(400, 200400, "name参数不能为空"))
		return
	}

	if req.Password == "" {
		context.Error(common.NewError(400, 200400, "password参数不能为空"))
		return
	}

	user := models.User{
		UserName: req.Name,
		Password: req.Password,
		Email:    req.Email,
		Devices:  []models.Device{},
	}

	err := user.Create()
	if err != nil {
		// 服务器错误处理
		context.Error(common.NewError(400, 200400, err.Error()))
	} else {
		success.Data = req.Name
		context.JSON(http.StatusOK, success)
	}
}

type myTag struct {
	Tag []string
}

// 创建用户接口

func UpdateUser(context *gin.Context) {
	var user *models.User
	var dbTag *models.Tag
	// 获取post方法的表单请求参数
	var tag myTag
	body, BErr := context.GetRawData()
	if BErr != nil {
		// 服务器错误处理
		context.Error(common.NewError(400, 200400, BErr.Error()))
		return
	}
	jsonErr := json.Unmarshal(body, &user)

	tagErr := json.Unmarshal(body, &tag)
	var tagsList []models.Tag
	// 有标签的情况下才设置
	if tagErr == nil {
		for _, v := range tag.Tag {
			tag := models.Tag{
				Name: v,
			}
			exist, res := dbTag.Exist((&tag))
			if exist {
				tagsList = append(tagsList, res)
			} else {
				tagsList = append(tagsList, tag)
			}
		}
		// 加载完之后呢
		user.Tags = tagsList
	}

	// 根据token获取用户的id
	userID, _ := context.Get("userID")

	user.ID = userID.(int)
	if tagErr != nil {
		// 服务器错误处理
		context.Error(common.NewError(400, 200400, jsonErr.Error()))
		return
	}
	if jsonErr != nil {
		// 服务器错误处理
		context.Error(common.NewError(400, 200400, jsonErr.Error()))
		return
	}
	err := user.Update()

	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	} else {
		res, _ := user.GetUserByID()
		context.JSON(200, gin.H{"mes": "success", "data": res})
		return
	}

}

// 根据名字获取用户信息
func GetUserByName(context *gin.Context) {
	name := context.Query("name")
	if name == "" {
		context.Error(common.NewError(400, 200400, "name参数不能为空"))
		return
	}
	user, err := MyUser.GetUserByName(name)

	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	}
	res := common.Success{
		Mes:  "success",
		Data: user,
	}
	context.JSON(200, res)
}

// 获取所有用户
func GetAllUser(context *gin.Context) {
	users, err := MyUser.GetAll()
	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	}
	res := common.Success{
		Mes:  "success",
		Data: users,
	}

	context.JSON(200, res)
}

// 根据id删除用户
func DeleteUserById(context *gin.Context) {
	id := context.Query("id")
	if id == "" {
		context.Error(common.NewError(400, 200400, "id参数不能为空"))
		return
	}
	newId, transfromErr := strconv.Atoi(id)
	if transfromErr != nil {
		context.Error(common.NewError(400, 200400, transfromErr.Error()))
		return
	}
	err := MyUser.DeleteUserById(newId)
	if err != nil {
		context.Error(common.NewError(400, 200400, err.Error()))
		return
	}
	res := common.Success{
		Mes:  "success",
		Data: "",
	}
	context.JSON(200, res)
}

// 根据token获取自己的信息
func GetUserInfoByTokenId(context *gin.Context) {

	strId, _ := context.Get("userID")
	userID, _ := strId.(int)
	res, _ := MyUser.GetUserByTokenID(userID)
	context.JSON(200, gin.H{"mes": "success", "data": res})
	return

}

func VertifyPw(context *gin.Context) {
	var namePd NamePW
	err := context.ShouldBindJSON(&namePd)
	if err != nil {
		context.Error(common.NewError(401, 200400, err.Error()))
	} else {
		_, vErr, userID := MyUser.VertifyNamePasswd(namePd.Password, namePd.Name)
		if vErr != nil {
			context.Error(common.NewError(400, 200400, "用户名或密码错误"))
			return
		} else {
			token, err := utils.GenerateToken(namePd.Name, userID)
			if err == nil {
				context.JSON(200, common.NewResp(200, 200, "认证成功", map[string]interface{}{
					"token": token,
					"id":    userID,
				}))

				return

			}
			context.JSON(401, gin.H{
				"msg":   "认证失败,生成秘钥失败",
				"token": "",
			})
			return
		}
	}

}
