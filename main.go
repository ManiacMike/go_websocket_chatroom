package main

import (
    "golang.org/x/net/websocket"
    "fmt"
    // "io"
    // "os"
    "log"
    "net/http"
    "strconv"
    "time"
)

type User struct {
	id string
	con *websocket.Conn
}

var Users []User

func generateId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func ChatWith(ws *websocket.Conn) {
    var err error
    uid := generateId()
    newUser := User{uid,ws}
    Users = append(Users,newUser)

    fmt.Println("connect current user num",len(Users))

    for {
        var reply string

        if err = websocket.Message.Receive(ws, &reply); err != nil {
            fmt.Println("Can't receive,user ",uid," lost connection")
            Users = removeUser(uid)
            break
        }

        fmt.Println("Received back from client: " + reply)

        //msg := "Received from " + ws.Request().Host + "  " + reply
        // msg := "welcome from websocket"
        // fmt.Println("Sending to client: " + msg)
        for _,user := range Users{
          if user.id == uid{
              continue
          }
          if err = websocket.Message.Send(user.con, reply); err != nil {
              fmt.Println("Can't send user ",user.id," lost connection")
              Users = removeUser(user.id)
              break
          }
        }
    }
}

func removeUser(uid string) []User{
	var find int
	flag := false
	for i,v:=range Users{
		if uid == v.id{
			find = i
			flag = true
			break
		}
	}
	if flag{
		newHay := append(Users[:find],Users[find+1:]...)
		return newHay
	}else{
		return Users
	}
}

func StaticServer(w http.ResponseWriter, req *http.Request) {
    http.ServeFile(w,req,"chat.html")
    // staticHandler := http.FileServer(http.Dir("./"))
    // staticHandler.ServeHTTP(w, req)
    return
}

func main() {

    http.Handle("/", websocket.Handler(ChatWith))
    http.HandleFunc("/chat", StaticServer)

    fmt.Println("listen on port 8001")
    fmt.Println("浏览器访问 http://127.0.0.1:8001/chat")

    if err := http.ListenAndServe(":8001", nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
