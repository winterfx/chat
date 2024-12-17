package rest

import (
	"chat/internal/repository"
	"chat/internal/services/chat"
	"chat/internal/services/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func Auth(c *gin.Context) {
	provider := c.Param("provider")
	a, err := user.NewOAuth(user.AuthProvider(provider), repository.GetDefaultUserRepo())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	url := a.GetAuthConsentUrl(c)
	c.Redirect(http.StatusFound, url)
}

func AuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	a, err := user.NewOAuth(user.AuthProvider(provider), repository.GetDefaultUserRepo())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	code := c.Query("code")
	token, err := a.ConnectIdp(c, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to retrieve user information"})
		return
	}
	u, err := a.RetrieveUserInfo(c, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jwtToken, err := user.GenerateToken(c, u.Id, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//redirect to home page
	SetTokenCookie(c, jwtToken)
	c.Redirect(http.StatusFound, fmt.Sprintf("http://localhost:5173/login"))
}

func CheckIsLogin(c *gin.Context) {
	token, err := c.Cookie("chat_token")
	if err == nil && token != "" {
		// 验证 token 并获取用户 ID
		uid, err := user.ValidateToken(c, token)
		u, err := repository.GetDefaultUserRepo().GetUserById(c, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err == nil {
			if user.IsTokenExpiringSoon(c, token) {
				// 生成新的 token
				newToken, err := user.GenerateToken(c, uid, 0)
				if err == nil {
					// 设置新的 token 到 cookie
					c.SetCookie("chat_token", newToken, 3600, "/", "localhost", false, true)
				}
			}

			c.JSON(http.StatusOK, &LoginResponse{Uid: uid, Name: u.Name})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"error": "User not logged in"})
}

func Login(c *gin.Context) {
	req := &LoginRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}
	uid, err := userService.Login(c, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jwtToken, err := user.GenerateToken(c, uid, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	SetTokenCookie(c, jwtToken)
	userState.RefreshUserState(c, uid, user.Online)
	c.JSON(http.StatusOK, &LoginResponse{Uid: uid})
}

func Register(c *gin.Context) {
	req := &RegisterRequest{}
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}
	uid, err := userService.RegisterUser(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jwtToken, err := user.GenerateToken(c, uid, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	SetTokenCookie(c, jwtToken)
	c.JSON(http.StatusOK, &RegisterResponse{Uid: uid})
}

func waitSubReady(c context.Context, uid string, conn *websocket.Conn) {
	log.Println("subscribe user channel")
	subReady := make(chan struct{})
	go chatService.SubscribeUserChanel(c, uid, chat.GetWebsocketMessageHandler(conn), subReady)
	<-subReady
	return
}

func waitConnectionClose(c context.Context, conn *websocket.Conn) {
	readClose := make(chan struct{})
	writeClose := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		read(c, conn, readClose, writeClose)
	}()
	go func() {
		defer wg.Done()
		write(c, conn, writeClose, readClose)
	}()
	wg.Wait()
}

func handleWs(c *gin.Context) {
	log.Println("handle ws")
	conn, err := upgradeToWs(c.Writer, c.Request)
	if err != nil {
		log.Printf("upgrade to ws error:%s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//subscribe to the chat service
	uid := GetUid(c)

	waitSubReady(c, uid, conn)
	log.Println("subscribed user channel finished")
	waitConnectionClose(c, conn)
	if err = chatService.UnsubscribeUserChannel(c, uid); err != nil {
		log.Printf("unsubscribe user channel error:%s\n", err)
	}
	if err = conn.Close(); err != nil {
		log.Printf("close ws connection error:%s\n", err)
	}
	log.Println("ws connection closed")
}

func write(c context.Context, conn *websocket.Conn, write, readClosed chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("write close")
		ticker.Stop()
		close(write)
	}()
	for {
		select {
		case <-readClosed:
			log.Println("read closed,so trigger write close")
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("write ping error:%s\n", err)
				return
			} else {
				log.Printf("ping sent\n")
			}
		}
	}
}

func read(c context.Context, conn *websocket.Conn, read, writeClosed chan struct{}) {
	defer func() {
		log.Println("read close")
		close(read)
	}()
	uid := GetUid(c)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		log.Printf("pong received %s\n", uid)
		err := userState.RefreshUserState(c, uid, user.Online)
		if err != nil {
			log.Printf("refresh user state error:%s\n", err)
		}
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		select {
		case <-writeClosed:
			log.Println("write closed,so trigger read close")
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				//closed
				log.Printf("read message error,channel closed:%s\n", err)
				if err = userState.RefreshUserState(c, uid, user.Offline); err != nil {
					log.Printf("refresh user state error:%s", err)
				}
				return
			}
			chatMessage := &chat.ChatMessage{}
			err = json.Unmarshal(message, chatMessage)
			if err != nil {
				fmt.Printf("unmarshal chat messsage error:%s\n", err)
			}
			err = chatService.PubMessage(c, chatMessage)
			if err != nil {
				fmt.Printf("pub msg err:%s\n", err)
			}
		}

	}
}

func RetrieveFriends(c *gin.Context) {
	uid := GetUid(c)
	fs, err := userService.GetUsersList(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	friends := make([]*ListFriendResponse, 0, len(fs))
	for _, f := range fs {
		s := userState.GetUserState(c, f.UserID)
		friends = append(friends, &ListFriendResponse{Uid: f.UserID, Name: f.Name, Status: s})
	}
	c.JSON(http.StatusOK, friends)
}
