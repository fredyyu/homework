package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/spf13/viper"
	"homework/database"
	"log"
	"net/http"
)

var client *linebot.Client

type User struct {
	UserId         string
	Message        string
	Name           string
	PictureURL     string
	ProfileMessage string
}

type SendForm struct {
	UserId  string `json:"userId" binding:required`
	Message string `json:"message" binding:required`
}

func InitRouter() *gin.Engine {

	lineClient, err := linebot.New(viper.GetString("LINE_CHANNEL_SECERT"), viper.GetString("LINE_CHANNEL_ACCESS"))

	if err != nil {
		log.Println(err.Error())
	}
	client = lineClient

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": http.StatusNotFound,
			"msg":    http.StatusText(http.StatusNotFound),
			"data":   make([]interface{}, 0),
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": http.StatusNotFound,
			"msg":    http.StatusText(http.StatusNotFound),
			"data":   make([]interface{}, 0),
		})
	})

	r.POST("/linehandler", lineHandler)
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "test")
	})

	r.POST("/send", send)
	r.GET("/user", user)

	return r
}

func lineHandler(c *gin.Context) {
	// recive request
	events, err := client.ParseRequest(c.Request)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(500)
		}

		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// reply message
				//get line webhook message all content
				// jsonEvent, _ := json.Marshal(event)
				// eventStr := string(jsonEvent)
				userId := event.Source.UserID
				res, err := client.GetProfile(userId).Do()
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				userInfo := User{userId, message.Text, res.DisplayName, res.PictureURL, res.StatusMessage}
				if event.Source.Type == "user" {
					if _, err = client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(res.DisplayName+":"+message.Text)).Do(); err != nil {
						log.Println(err.Error())
					}
				} else if event.Source.Type == "group" {
					if _, err = client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(res.DisplayName+":"+message.Text)).Do(); err != nil {
						log.Println(err.Error())
					}
				}

				err = database.InsertUserInfo("db", "user_info", userInfo)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}
	}
}

func send(c *gin.Context) {
	form := new(SendForm)
	if err, msg := bindandvalidate(c, form); err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status": 100,
			"msg":    msg,
			"data":   make([]interface{}, 0),
		})
		return
	}
	fmt.Println(form.UserId)
	fmt.Println(form.Message)
	_, err := client.PushMessage(form.UserId, linebot.NewTextMessage(form.Message)).Do()
	if err != nil {
		log.Printf("Send message to %s fail", form.UserId)
		c.JSON(http.StatusOK, map[string]interface{}{
			"status": 100,
			"msg":    err.Error(),
			"data":   make([]interface{}, 0),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"msg":    "success",
		"data":   make([]interface{}, 0),
	})
	return
}

func user(c *gin.Context) {

}

func bindandvalidate(c *gin.Context, form interface{}) (err error, msg string) {
	if err = c.ShouldBindJSON(form); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, fieldError := range errs {
			msg = fieldError.Error()
			return
		}
	}
	return
}
