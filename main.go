package main

import (
    "golang.org/x/net/websocket"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
    "encoding/json"
    "github.com/bitly/go-simplejson" // for json get
)

type User struct {
	id string
	con *websocket.Conn
}

var Users []User

// type Receive struct{
//     content   string    `json:"content"`
//     uname  string   `json:"uname"`
// }

type Reply struct{
    //id string
    Uname string  `json:"uname"`
    Content string  `json:"content"`
    Time int64  `json:"time"`
}

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
        var receiveMsg string

        if err = websocket.Message.Receive(ws, &receiveMsg); err != nil {
            fmt.Println("Can't receive,user ",uid," lost connection")
            Users = removeUser(uid)
            break
        }
        receive, err := simplejson.NewJson([]byte(receiveMsg))
        if err != nil {
            panic(err.Error())
        }
        var receiveNodes = make(map[string]interface{})
        receiveNodes, _ = receive.Map()
        fmt.Println("Received back from client: " , receiveNodes)

        //msg := "Received from " + ws.Request().Host + "  " + reply
        reply := Reply{Uname:receiveNodes["uname"].(string),Content:receiveNodes["content"].(string),Time:time.Now().Unix()}
        //fmt.Println(reply)
        replyBody, err := json.Marshal(reply)
        if err != nil {
            panic(err.Error())
        }
        //fmt.Println(replyBody)
        replyBodyStr := string(replyBody);
        //fmt.Println(replyBodyStr)
        for _,user := range Users{
          // if user.id == uid{
          //     continue
          // }
          if err = websocket.Message.Send(user.con, replyBodyStr); err != nil {
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
