package ws

import (
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
var clientMap map[string]*Node = make(map[string]*Node, 0)

// 消息类型
const GET_MSG = "GET_MSG"
const LOGIN_SUCCESS = "LOGIN_SUCCESS"
const HEART_BEAT = "HEART_BEAT"
const SENDMSG = "SEND_MSG"
const CLOSECLIENT = "CLOSE_CLIENT"
const UPDATE_USER_LIST = "UPDATE_USER_LIST"

// 读写锁
var rwLocker sync.RWMutex

// 消息类型
type ContactMes struct {
	MesType    string `json:"type"`
	FromUserID int    `json:"userID"`
	Mes        string `json:"mes"`
	ToUserID   string `json:"toUserID"`
}

// 消息类型
type HeartBeat struct {
	MesType string `json:"type"`
	Mes     string `json:"mes"`
}

// Node 当前用户节点 userId和Node的映射关系
type Node struct {
	ClientID string `json:"name"`
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
	var client models.WsClient
	client.User = claims.UserID
	fmt.Println("=========")
	fmt.Println(claims.UserID)
	client.GetClientByUser()
	node := &Node{
		ClientID:     client.WebsocketID,
		WsClientInfo: client,
		Conn:         conn,
		DataQueue:    make(chan interface{}, 50),
	}
	// 映射关系的绑定
	clientMap[client.WebsocketID] = node
	// 广播更新用户列表
	broadcast()
	rwLocker.Unlock()
	sendMsg(client.WebsocketID, map[string]interface{}{"msg": "init success", "user_id": client.WebsocketID, "type": LOGIN_SUCCESS, "users": clientMap})
	// 发送数据给客户端
	go sendProc(node)
	// 接收消息
	go listenFromClient(node)
}

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
		// 发送消息处理
		case SENDMSG:
			// 发送消息给某个用户,判断该用户是否还在
			tmpNode, Nerr := clientMap[resMes.ToUserID]
			// 如果用户已经
			if !Nerr {
				// 给发送者回一条消息
				node.DataQueue <- map[string]interface{}{"type": GET_MSG, "to": resMes.ToUserID, "from": 0, "msg": "该用户已下线"}

			} else {
				tmpNode.DataQueue <- map[string]interface{}{"type": GET_MSG, "to": resMes.ToUserID, "from": resMes.FromUserID, "msg": resMes.Mes}
			}

		}
	}
}

// 将数据推送到管道中
func sendMsg(clientID string, message interface{}) {
	rwLocker.RLock()
	node, isOk := clientMap[clientID]
	fmt.Println(node)
	rwLocker.RUnlock()
	if isOk {
		node.DataQueue <- message
	}
}

// 从管道中获取数据发送出去
// 心跳保活机制
func sendProc(node *Node) {
	timer := time.NewTicker(5 * time.Second) // 5s后触发
	for {
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
			fmt.Println("超时了")
			timer.Stop()
			// 超时了就把这个链接关闭，然后去掉这个map
			node.Conn.Close()
			fmt.Printf("%s已关闭", node.ClientID)
			delete(clientMap, node.ClientID)
			// 广而告之，有用户退出了
			broadcast()
			goto EXIT
		}
	}
EXIT:
}

// 广播更新列表
func broadcast() {
	for _, v := range clientMap {
		v.DataQueue <- map[string]interface{}{"type": UPDATE_USER_LIST, "users": clientMap}
	}
}
