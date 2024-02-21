package ws

import (
	"GetHotWord/common"
	"GetHotWord/interval/api/models"
	"GetHotWord/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 映射关系表
var clientMap map[int]*Node = make(map[int]*Node, 0)

// 消息类型
const GET_MSG = "GET_MSG"
const LOGIN_SUCCESS = "LOGIN_SUCCESS"
const HEART_BEAT = "HEART_BEAT"
const SENDMSG = "SEND_MSG"
const CLOSECLIENT = "CLOSE_CLIENT"
const UPDATE_USER_LIST = "UPDATE_USER_LIST"
const RTC_SDP = "RTC_SDP"

// 读写锁
var rwLocker sync.RWMutex

type SDPData struct {
	Sdp  string `json:"sdp"`
	Type string `json:"type"`
}

type SDP struct {
	From int     `json:"from"`
	To   int     `json:"to"`
	Data SDPData `json:"sdp"`
}

// 消息类型
type ContactMes struct {
	MesType string         `json:"type"`
	Data    models.Message `json:"data"`
	SDP     SDP            `json:"sdpData"`
}

// 消息类型
type HeartBeat struct {
	MesType string `json:"type"`
	Mes     string `json:"mes"`
}

// Node 当前用户节点 userId和Node的映射关系
type Node struct {
	UserID int `json:"userID"`
	// 这个是保留该node的wsid
	WsClientInfo models.WsClient `json:"wsClientInfo"`
	//这个是维护链接
	Conn *websocket.Conn `json:"-"`
	// 这个是消息队列
	DataQueue chan interface{} `json:"-"`
	// 群组的消息分发
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 根据判断token的方法来鉴权,如果没token就返回false
		return true
	},
}

func Chat(ctx *gin.Context) {
	// 升级协议以后原来的请求头会消失，所以要在query里用来获取
	token := ctx.Query("token")
	_, claims, _ := utils.ParseToken(token)

	// 在响应头上添加Sec-Websocket-Protocol,
	upGrader.Subprotocols = []string{ctx.GetHeader("Sec-Websocket-Protocol")}
	//升级get请求为webSocket协议
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		return
	}

	// 绑定到当前节点
	rwLocker.Lock()
	// var client models.WsClient
	// client.User = claims.UserID
	// client.GetClientByUser()
	if err != nil {
		ctx.Error(common.NewError(500, 200500, "系统出错了"))
		return
	}

	_, exists := clientMap[claims.UserID]

	// 判断是不是首次登录，不是首次的话之前应该是有相应的记录的
	if exists {
		// 更新的东西  1、状态  2、双方的连接
		clientMap[claims.UserID].WsClientInfo.Online = 1
		clientMap[claims.UserID].Conn = conn
		clientMap[claims.UserID].DataQueue = make(chan interface{}, 50)
	} else {
		var client models.WsClient
		client.User = claims.UserID
		client.GetClientByUser()
		node := &Node{
			UserID:       claims.UserID,
			WsClientInfo: client,
			Conn:         conn,
			DataQueue:    make(chan interface{}, 50),
		}
		// 映射关系的绑定
		clientMap[claims.UserID] = node
	}
	// 广播更新用户列表
	broadcast()
	rwLocker.Unlock()
	sendMsg(claims.UserID, map[string]interface{}{
		"msg":      "init success",
		"userID":   claims.UserID,
		"clientID": clientMap[claims.UserID].WsClientInfo.WebsocketID,
		"type":     LOGIN_SUCCESS,
		"userInfo": clientMap[claims.UserID].WsClientInfo.UserInfo,
		"users":    clientMap,
	})
	// 发送数据给客户端
	go sendProc(clientMap[claims.UserID])
	// 接收消息
	go listenFromClient(clientMap[claims.UserID])
}

/*
*

	监听客户端的socket是否有发送消息过来，并处理发送的消息

*
*/
func listenFromClient(node *Node) {
	for {
		_, mes, err := node.Conn.ReadMessage()
		if err != nil {
			node.Conn.WriteJSON(map[string]interface{}{"type": "error", "msg": err})
			return
		}
		var resMes ContactMes
		json.Unmarshal(mes, &resMes)
		switch resMes.MesType {
		// 心跳包处理
		case HEART_BEAT:
			// 给本人发送pong消息
			node.DataQueue <- (map[string]interface{}{"type": HEART_BEAT, "msg": "pong"})

			// sdp信息交换，由于webrtc
		case RTC_SDP:
			// 给远端发送消息
			clientMap[resMes.SDP.To].DataQueue <- (map[string]interface{}{
				"from": resMes.SDP.From,
				"name": clientMap[resMes.SDP.From].WsClientInfo.UserInfo.UserName,
				"type": RTC_SDP,
				"msg":  resMes.SDP.Data,
			})
			// 发送消息处理

		case SENDMSG:
			// 发送消息给某个用户,判断该用户是否还在
			tmpNode, Nerr := clientMap[resMes.Data.To]
			message := models.Message(resMes.Data)
			res, err := message.Create()
			if err != nil {
				// 给发送者回一条消息
				node.DataQueue <- map[string]interface{}{
					"type":    GET_MSG,
					"to":      resMes.Data.To,
					"message": res,
					"from":    0,
					"msg":     "数据保存失败",
				}
			}
			// 如果用户已经
			if !Nerr {
				// 给发送者回一条消息
				node.DataQueue <- map[string]interface{}{
					"type": GET_MSG,
					"to":   resMes.Data.To,
					"from": 0,
					"msg":  "该用户已下线",
				}
			} else {
				node.DataQueue <- map[string]interface{}{
					"type":    GET_MSG,
					"to":      resMes.Data.To,
					"from":    resMes.Data.From,
					"message": res,
					"msg":     "success",
				}
				tmpNode.DataQueue <- map[string]interface{}{
					"type":    GET_MSG,
					"to":      resMes.Data.To,
					"from":    resMes.Data.From,
					"message": res,
					"msg":     resMes.Data.Content,
				}
			}

		}
	}
}

// 将数据推送到管道中
func sendMsg(userID int, message interface{}) {
	rwLocker.RLock()
	node, isOk := clientMap[userID]
	fmt.Println(node)
	rwLocker.RUnlock()
	if isOk {
		node.DataQueue <- message
	}
}

// 给该用户的数据都都会发到所属管道里，管道中获取数据发送该用户
// 心跳保活机制，如果是时钟先到就结束了该函数
func sendProc(node *Node) {
	timer := time.NewTicker(5 * time.Second) // 5s后触发
	// 加锁保护连接关闭操作
	rwLocker.Lock()
	conn := node.Conn
	rwLocker.Unlock()
	// 无限循环
EXIT:
	for {
		// 看看是时钟先到还是心跳先到
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteJSON(data)
			if err != nil {
				fmt.Println(err)
			} else {
				timer.Stop()
				timer.Reset(5 * time.Second)
			}

		// 判断有没有超时保活
		case <-timer.C:
			timer.Stop()
			// 超时了就把这个链接关闭，然后置为下线,
			// node.Conn == conn这个是为了防止用户刷新了页面,导致conn已经被更换
			rwLocker.Lock()
			if node.Conn == conn {
				fmt.Printf("%d已关闭", node.UserID)
				conn.Close()
				clientMap[node.UserID].WsClientInfo.Online = 0
				broadcast()
			}
			rwLocker.Unlock()
			// 退出这个循环，因为用户都断线了
			break EXIT
		}
	}

}

// 广播更新列表
func broadcast() {
	for _, v := range clientMap {
		// 在线的才发送
		if v.WsClientInfo.Online == 1 {
			v.DataQueue <- map[string]interface{}{
				"type":  UPDATE_USER_LIST,
				"users": clientMap,
			}
		}

	}
}
