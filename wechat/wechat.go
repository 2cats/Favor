package wechat

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"favor/model"

	"github.com/chanxuehong/wechat/corp"
	"github.com/chanxuehong/wechat/corp/menu"
	"github.com/chanxuehong/wechat/corp/message/request"
	"github.com/chanxuehong/wechat/corp/message/response"
	"github.com/chanxuehong/wechat/util"
)

var corpClient *corp.Client

const (
	NEWEST = "NEWEST"
)

func TextMessageHandler(w http.ResponseWriter, r *corp.Request) {
	text := request.GetText(r.MixedMsg)
	ret := "err"
	err := models.InsertMsg(text.Content, text.FromUserName)
	if err == nil {
		ret = "success"
	}
	// content := strings.TrimSpace(text.Content)
	// id, err := strconv.Atoi(content)
	// ret := "err"
	// if err == nil {
	// 	msg, err := models.SelectIdMsg(id)
	// 	if err == nil {
	// 		ret = msg.Content
	// 	}
	// }
	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, ret)
	corp.WriteResponse(w, r, resp)
}
func url2File(url string, filename string) error {
	response, e := http.Get(url)
	if e != nil {
		return e
	}
	defer response.Body.Close()
	// prefix := strconv.FormatInt(time.Now().Unix(), 10) + "_"

	//open a file for writing
	file, err := os.Create(filename)
	if err != nil {
		return e
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	file.Close()
	return err
}
func ImageMessageHandler(w http.ResponseWriter, r *corp.Request) {
	image := request.GetImage(r.MixedMsg)
	ret := "err"
	path := models.UPLOAD_DIR + models.GetFilePrefix() + image.FromUserName + ".png"
	err := url2File(image.PicURL, path)
	if err == nil {
		err = models.InsertMsg(fmt.Sprintf(`<img src="%s"></img>`, path), image.FromUserName)
		if err == nil {
			ret = "success"
		}

	}
	resp := response.NewText(image.FromUserName, image.ToUserName, image.CreateTime, ret)
	corp.WriteResponse(w, r, resp)
}
func ClickEventHandler(w http.ResponseWriter, r *corp.Request) {
	ret := "err"
	switch r.MixedMsg.EventKey {
	case NEWEST:
		msg, err := models.SelectNewestMsg()
		if err == nil {
			ret = msg.Content
		}
	}
	resp := response.NewText(r.MixedMsg.FromUserName, r.MixedMsg.ToUserName, r.MixedMsg.CreateTime, ret)
	corp.WriteResponse(w, r, resp)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}
func WechatInit() http.Handler {
	aesKey, err := util.AESKeyDecode("Cq6kwC4ZXIeOARp12yGygfk5JXcjr4GgWnnUYoXAsUC")
	if err != nil {
		panic(err)
	}
	messageServeMux := corp.NewMessageServeMux()
	messageServeMux.MessageHandleFunc(request.MsgTypeText, TextMessageHandler)
	messageServeMux.MessageHandleFunc(request.MsgTypeImage, ImageMessageHandler)

	messageServeMux.EventHandleFunc("click", ClickEventHandler)

	agentServer := corp.NewDefaultAgentServer("wx26ec68469d8c2881", 5 /* agentId */, "amber", aesKey, messageServeMux)
	agentServerFrontend := corp.NewAgentServerFrontend(agentServer, corp.ErrorHandlerFunc(ErrorHandler), nil)

	var AccessTokenServer = corp.NewDefaultAccessTokenServer("wx26ec68469d8c2881", "cqsRM4HAJajHQniwKgyvosxQu_bYhaj1t9czyFxj5oGO0Eg9pKCfJcF9SiyHHyjk", nil)
	corpClient = corp.NewClient(AccessTokenServer, nil)

	var mn menu.Menu
	mn.Buttons = make([]menu.Button, 1)
	mn.Buttons[0].SetAsClickButton("查看最新", NEWEST)

	menuClient := (*menu.Client)(corpClient)

	if err := menuClient.CreateMenu(5 /* agentId */, mn); err != nil {
		fmt.Println(err)
		return nil
	}
	return agentServerFrontend
	// http.Handle("/agent", agentServerFrontend)
	// log.Fatal(http.ListenAndServe(":8889", nil))

}
