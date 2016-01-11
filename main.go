package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

var CurrentUsers *UserList //在线用户列表

type MessageReply struct {
	Type    string `json:"type"`
	Uname   string `json:"uname"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type UidCookieReply struct {
	Type string `json:"type"`
	Uid  string `json:"uid"`
}

type UserCountChangeReply struct {
	Type      string `json:"type"`
	UserCount int    `json:"user_count"`
}

func ChatServer(ws *websocket.Conn) {
	var err error
	var uid string
	if nil == CurrentUsers {
		CurrentUsers = new(UserList)
	}
	if uidCookie, err := ws.Request().Cookie("uid"); err != nil {
		fmt.Println("visitor is unknown")
		uid = CurrentUsers.New(ws)
	} else {
		uid := uidCookie.Value
		fmt.Println("visitor ", uid)
		userExist, index := CurrentUsers.Exist(uid)
		if userExist == true {
			CurrentUsers.ChangeConn(index, ws)
		} else {
			fmt.Println("visitor uid is outdate")
			uid = CurrentUsers.New(ws) //cookie中的uid不存在
		}
	}
	go PushUserCount()

	for {
		var receiveMsg string

		if err = websocket.Message.Receive(ws, &receiveMsg); err != nil {
			fmt.Println("Can't receive,user ", uid, " lost connection")
			CurrentUsers.Remove(uid)
			break
		}

		receiveNodes := JsonStrToStruct(receiveMsg)
		fmt.Println("Received back from client: ", receiveNodes)
		reply := MessageReply{Type: "message", Uname: receiveNodes["uname"].(string), Content: receiveNodes["content"].(string), Time: time.Now().Unix()}
		replyBody, err := json.Marshal(reply)
		if err != nil {
			panic(err.Error())
		}
		replyBodyStr := string(replyBody)
		go Broadcast(replyBodyStr)
	}
}

func PushUserCount() {
	userCount := UserCountChangeReply{"user_count", len(*CurrentUsers)}
	replyBody, err := json.Marshal(userCount)
	if err != nil {
		panic(err.Error())
	}
	replyBodyStr := string(replyBody)
	Broadcast(replyBodyStr)
}

func Broadcast(replyBodyStr string) error {
	fmt.Println("current user", len(*CurrentUsers))
	for _, user := range *CurrentUsers {
		if err := websocket.Message.Send(user.con, replyBodyStr); err != nil {
			fmt.Println("Can't send user ", user.uid, " lost connection")
			CurrentUsers.Remove(user.uid)
			break
		}
	}
	return nil
}

func StaticServer(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "chat.html")
	// staticHandler := http.FileServer(http.Dir("./"))
	// staticHandler.ServeHTTP(w, req)
	return
}

func main() {

	http.Handle("/", websocket.Handler(ChatServer))
	http.HandleFunc("/chat", StaticServer)

	fmt.Println("listen on port 8001")
	fmt.Println("浏览器访问 http://yourhost:8001/chat")

	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
