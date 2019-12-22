package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Code    int    `json:"code"`
	Details string `json:"details"`
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var interval time.Duration
var ORIGIN string
var PORT int
var HOST string
//var Don_update bool

type Host struct {
	Hostname string `form:"hostname" json:"hostname" binding:"required"`
	Port     string `form:"port,default=22" json:"port" `
	User     string `form:"user,default=root" json:"user"`
	Password string `form:"pd" json:"password"` //base64 encoded password
	Cols     int    `form:"cols,default=120" json:"cols"`
	Rows     int    `form:"rows,default=32" json:"rows"`
	PriKey   string `form:"pk" json:"pri_key"` //base64 private key
}

func main() {
	Init()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	e := gin.Default()
	e.Use(cors.New(setOrigins()))
	e.GET("/ws", WsSsh)
	e.Run(fmt.Sprintf("%s:%d",HOST, PORT))

}
func Init() {
	flag.DurationVar(&interval, "interval", 20*time.Second, "set ping pong frequency")
	flag.StringVar(&ORIGIN, "origin", "*", "set origins ,like \"http://127.0.0.1:8080,http://localhost:8080\"")
	flag.IntVar(&PORT,"port",9018,"set port")
	flag.StringVar(&HOST,"host","","set bind host (default all)")
	/*var v = flag.Bool("V", false, "show version")
	if *v {
		showVersion()
	}*/
	flag.Parse()
	log.Printf("interval is %s", interval)
}
func setOrigins() cors.Config {
	config := cors.DefaultConfig()
	if ORIGIN !="*" {
		config.AllowOrigins=parseOrigins(ORIGIN)
		log.Println(parseOrigins(ORIGIN))
	} else {
		config.AllowAllOrigins = true
	}
	return config
}
func parseOrigins(o string) ([]string) {
	oList:=strings.Split(o,",")
	return oList
}
func wshandleError(ctx *gin.Context, err error) (bool) {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Message{http.StatusInternalServerError,
			fmt.Sprintf("%s", err)})
	}
	return (err != nil)
}
func WsSsh(c *gin.Context) {
	var host Host
	err := c.ShouldBindQuery(&host)
	if wshandleError(c, err) {
		return
	}
	if host.Password != "" {
		pD, _ := base64.StdEncoding.DecodeString(host.Password)
		host.Password = string(pD)
	}
	var key *Key
	key = new(Key)
	if host.PriKey != "" {
		sD, _ := base64.URLEncoding.DecodeString(host.PriKey)
		key.PrivateKey = string(sD)
	} else {
		key = nil
	}
	client, err := NewSshClient(host.User, host.Hostname, host.Port, host.Password, key)

	if wshandleError(c, err) {
		return
	}
	defer client.Close()

	ssConn, err := NewSshConn(host.Cols, host.Rows, client)
	if wshandleError(c, err) {
		return
	}
	defer ssConn.Close()
	// after configure, the WebSocket is ok.
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if wshandleError(c, err) {
		return
	}
	defer wsConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan
}
