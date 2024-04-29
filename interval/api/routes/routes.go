package routes

import (
	// cors2 "GetHotWord/config/cors"
	"github.com/wuyz-harder/notebook-backend/interval/api/controller/comment"
	"github.com/wuyz-harder/notebook-backend/interval/api/controller/device"
	"github.com/wuyz-harder/notebook-backend/interval/api/controller/file"
	"github.com/wuyz-harder/notebook-backend/interval/api/controller/message"
	"github.com/wuyz-harder/notebook-backend/interval/api/controller/note"
	"github.com/wuyz-harder/notebook-backend/interval/api/controller/star"
	"github.com/wuyz-harder/notebook-backend/interval/api/controller/tag"
	"github.com/wuyz-harder/notebook-backend/interval/api/controller/user"
	"github.com/wuyz-harder/notebook-backend/interval/ws"

	"strings"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	// 跨域
	// r.Use(cors.New(cors2.GetCors()))

	strings.Contains("222", "22")
	v1 := r.Group("/v1/api")
	// 设备
	v1.Handle("GET", "/device-user", device.GetDeviceByUserId)
	// 获取某个设备的信息
	v1.Handle("GET", "/device", device.GetDeviceById)
	// 修改某个设备的信息
	v1.Handle("POST", "/device/update", device.UpdateDevice)
	// 获取所有设备的信息
	v1.Handle("GET", "/devices", device.GetDevices)
	// 删除某个设备根据id
	v1.Handle("DELETE", "/device", device.DeleteDeviceById)
	// 添加设备
	v1.Handle("POST", "/device", device.AddDevice)
	// 上传文件
	v1.Handle("POST", "/upload", file.UploadFile)
	// 上传头像
	v1.Handle("POST", "/avatar", file.UploadAvatar)

	// 用户
	// 复数代表获取所有用户
	v1.Handle("GET", "/users", user.GetAllUser)
	// 复数代表获取所有用户
	v1.Handle("POST", "/user/update", user.UpdateUser)
	// 根据用户名获取某个用户的详细信息
	v1.Handle("GET", "/user", user.GetUserByName)
	// 根据token里的id获取本人的详细信息
	v1.Handle("GET", "/myselfInfo", user.GetUserInfoByTokenId)
	// 删除某个用户
	v1.Handle("DELETE", "/user", user.DeleteUserById)
	// 创建用户
	v1.Handle("POST", "/user/register", user.CreateUser)
	// 用户登录
	v1.Handle("POST", "/user/login", user.VertifyPw)

	// 标签
	// 复数代表获取所有
	v1.Handle("GET", "/tags", tag.GetAll)
	v1.Handle("POST", "/tags", tag.Add)
	v1.Handle("DELETE", "/tags", tag.Delete)

	// 聊天
	v1.Handle("GET", "/chat", ws.Chat)

	//笔记
	v1.Handle("POST", "/note", note.CreateNote)
	v1.Handle("POST", "/note/share", note.ShareNote)
	// 根据share的URL去获取笔记
	v1.Handle("GET", "/note/share", note.GetShareNote)
	v1.Handle("GET", "/note", note.AllNote)

	v1.Handle("DELETE", "/note/:id", note.DeleteNote)

	// 评论
	v1.Handle("POST", "/comment", comment.CreateComment)
	// 获取某个笔记下的所有评论
	v1.Handle("GET", "/comment", comment.GetComments)

	// 收藏
	v1.Handle("POST", "/star", star.CreateStar)
	v1.Handle("GET", "/star", star.GetStars)

	// 消息记录
	v1.Handle("GET", "/message", message.GetMessages)

}
