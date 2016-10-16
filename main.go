package main

import (
	"fmt"
	"favor/model"
	"favor/wechat"
	"strconv"

	"github.com/kataras/iris"
	"github.com/valyala/fasthttp"
)

type CommonReponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

var (
	EmptySuccessResponse = CommonReponse{"SUCCESS", nil}
)

func main() {
	models.DBInit()
	iris.Config.MaxRequestBodySize = 1024 * 1024 * 1024 //1G
	iris.Config.IsDevelopment = true
	agentHandler := wechat.WechatInit()
	iris.Handle("", "/agent", iris.ToHandler(agentHandler))
	iris.Get("/", func(ctx *iris.Context) {
		ctx.MustRender("index.html", struct {
			Msgs []string
		}{
			Msgs: []string{"aaa", "bbb"},
		})
	})
	iris.Static("/public", "./public/", 1)
	iris.Get("/msg", func(ctx *iris.Context) {
		pgStr := ctx.FormValueString("pagesize")
		nthStr := ctx.FormValueString("nth")
		filter := ctx.FormValueString("filter")

		if pgStr == "" && nthStr == "" {
			resp := EmptySuccessResponse
			count, err := models.GetMsgCount(filter)
			if err != nil {
				return
			}
			resp.Data = count
			ctx.JSON(200, resp)
			return
		}
		pagesize, err := strconv.Atoi(pgStr)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		nth, err := strconv.Atoi(nthStr)
		if err != nil {
			return
		}

		msg, err := models.SelectPageMsg(filter, pagesize, nth)
		if err != nil {
			return
		}
		resp := EmptySuccessResponse
		resp.Data = msg
		ctx.JSON(200, resp)
	})
	iris.Get("/newest", func(ctx *iris.Context) {
		msg, err := models.SelectNewestMsg()
		if err != nil {
			ctx.WriteString(err.Error())
		} else {
			ctx.HTML(200, msg.Content)

		}
	})
	iris.Post("/msg", func(ctx *iris.Context) {
		msg := struct {
			Op   string `json:"op"`
			Data string `json:"data"`
		}{}

		if err := ctx.ReadJSON(&msg); err != nil {
			return
		}

		switch msg.Op {
		case "INSERT":
			if err := models.InsertMsg(msg.Data, models.AdminUser); err != nil {
				return
			}
			ctx.JSON(200, EmptySuccessResponse)
		case "DELETE":

			id, err := strconv.ParseInt(msg.Data, 10, 0)
			if err != nil {
				return
			}
			if err := models.DeleteMsg(id); err != nil {
				return
			}
			fmt.Printf("\n %d DELETED \n", id)
			ctx.JSON(200, EmptySuccessResponse)
		}

	})
	iris.Post("/postfiles", func(ctx *iris.Context) {
		response := CommonReponse{}
		files := make(map[string]interface{})
		key := 0
		for {
			header, err := ctx.FormFile(strconv.Itoa(key))
			key++
			if err != nil {
				response = EmptySuccessResponse
				response.Data = files
				ctx.JSON(200, response)
				return
			}

			path := models.UPLOAD_DIR + models.GetFilePrefix() + header.Filename
			fasthttp.SaveMultipartFile(header, path)
			files[header.Filename] = path
		}
	})

	iris.Listen(":8888")
}
