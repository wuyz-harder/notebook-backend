package routes

import (
	// cors2 "GetHotWord/config/cors"
	"GetHotWord/interval/api/controller/comment"
	"GetHotWord/interval/api/controller/device"
	"GetHotWord/interval/api/controller/file"
	"GetHotWord/interval/api/controller/note"
	"GetHotWord/interval/api/controller/star"
	"GetHotWord/interval/api/controller/tag"
	"GetHotWord/interval/api/controller/user"
	"GetHotWord/interval/ws"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	// 跨域
	// r.Use(cors.New(cors2.GetCors()))

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
	v1.Handle("POST", "/user", user.CreateUser)
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

}
