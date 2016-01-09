package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "encoding/json"
)

type User struct {
	uid string
	con *websocket.Conn
}

type UserList []User

func (users *UserList) New(ws *websocket.Conn) string{
  uid := GenerateId()
  (*users) = append(*users,User{uid,ws})
  fmt.Println(*users)
  fmt.Println("New user connect current user num",len(*users))
  reply := UidCookieReply{Type:"session",Uid:uid}
  replyBody, err := json.Marshal(reply)
  if err != nil {
      panic(err.Error())
  }
  replyBodyStr := string(replyBody)
  if err := websocket.Message.Send(ws, replyBodyStr); err != nil {
      fmt.Println("Can't send user ",uid," lost connection")
      users.Remove(uid)
  }
  return uid
}

func (users *UserList)Remove(uid string){
  flag,find := users.Exist(uid)
	if flag == true{
		(*users) = append((*users)[:find],(*users)[find+1:]...)
    PushUserCount()
	}
}

func (users *UserList)ChangeConn(index int,con *websocket.Conn){
  fmt.Println("visitor exist change connection")
  curUser := (*users)[index]
  curUser.con.Close()
  (*users)[index].con = con
}

func (users *UserList)Exist(uid string) (bool,int){
  var find int
  flag := false
  for i,v:=range *users{
    if uid == v.uid{
      find = i
      flag = true
      break
    }
  }
  return flag,find
}
