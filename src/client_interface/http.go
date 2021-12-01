package client_interface

import (
	"dejavuDB/src/auth"
	"dejavuDB/src/config"
	"dejavuDB/src/datastore"
	"dejavuDB/src/javascriptAPI"
	"dejavuDB/src/message"
	"dejavuDB/src/network"
	"dejavuDB/src/user"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var websocket_UpGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Http_client() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	authed := router.Group("/", BasicAuth)
	authed.GET("/storage", handleGet)
	authed.POST("/storage")
	authed.DELETE("/storage")
	authed.PUT("/storage")
	authed.PATCH("/storage")
	authed.HEAD("/storage")
	authed.Any("/adm")
	authed.Any("/sql")
	authed.Any("/js")
	authed.Any("/websocket", handleWebsocket)
	authed.Any("/tcp")
	//router.RunTLS()
	router.Run("127.0.0.1:" + config.Client_port)
}

func BasicAuth(c *gin.Context) {
	realm := "Basic realm=" + strconv.Quote("Authorization Required")
	cred := c.Request.Header.Get("Authorization")
	if cred == "" || len(cred) < 7 || cred[:7] != "Basic " {
		c.Header("WWW-Authenticate", realm)
		c.AbortWithStatus(401)
		return
	}
	cred = strings.Replace(cred, "Basic ", "", 1)
	b, err := base64.StdEncoding.DecodeString(cred)
	if err != nil {
		c.Header("WWW-Authenticate", realm)
		c.AbortWithStatus(401)
		return
	}
	u := strings.SplitN(string(b), ":", 1)
	if len(u) != 2 {
		c.Header("WWW-Authenticate", realm)
		c.AbortWithStatus(401)
		return
	}
	_, ok := user.Login(u[0], u[1])
	if !ok {
		c.Header("WWW-Authenticate", realm)
		c.AbortWithStatus(401)
		return
	}
	log.Println("User logged in: " + u[0])
	c.Set(gin.AuthUserKey, u[0])
}

func handleGet(c *gin.Context) {
	// get user, it was set by the BasicAuth middleware
	username := c.MustGet(gin.AuthUserKey).(string)

	key, ok := c.Get("key")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("key not provided"))
		return
	}
	if !auth.HasPermission(username, key.(string)) {
		c.AbortWithStatus(401)
		return
	}
	switch config.Role {
	case "groupLeader":
	case "groupMember":
	case "standalone":
		c.Writer.Write(datastore.Get(key.(string)).ToBytes())
	}
}

func handleWebsocket(c *gin.Context) {
	username := c.MustGet(gin.AuthUserKey).(string)
	ws, err := websocket_UpGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	for {
		m, msg, err := ws.ReadMessage()
		if err != nil {
			ws.WriteMessage(
				websocket.BinaryMessage,
				message.NewErrorClientMessage(err).ToBytes())
			return
		}
		if m != websocket.BinaryMessage {
			continue
		}
		handleClientMessage(msg, username)
	}
}

func handleTCP(c *gin.Context) {
	username := c.MustGet(gin.AuthUserKey).(string)

	h, ok := c.Writer.(http.Hijacker)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	netconn, conbuf, err := h.Hijack()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	netconn.SetDeadline(time.Time{})
	if conbuf.Reader.Buffered() != 0 {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	p := []byte{}
	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: tcp\r\nConnection: Upgrade\r\n"...)

	netconn.Write(p)

	conn := network.NewConn(netconn)
	defer conn.Conn.Close()

	err = conn.ClientHandshake()
	if err != nil {
		return
	}
	for {
		b, err := conn.ReadMesaage()
		if err != nil {
			conn.Write([]byte(err.Error()))
			return
		}
		re, err := handleClientMessage(b, username)
		if err != nil {
			conn.Write(message.NewErrorClientMessage(err).ToBytes())
		} else {
			conn.Write(re)
		}
	}
}

func handleClientMessage(b []byte, username string) ([]byte, error) {
	msg := message.ClientMessageFromBytes(b)

	switch msg.Type {
	case message.JsMessageType:
		return javascriptAPI.JavascriptRun(
			string(msg.Content),
			javascriptAPI.JsOptions{
				UserName: username,
			})
	}
	return nil, nil
}
