package main

import (
	"favor/model"
	"path/filepath"
	"strconv"

	"os"

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
	if len(os.Args) > 1 {
		models.ConfigFileName = os.Args[1]
	} else {
		models.ConfigFileName = models.DefaultConfFile
	}
	models.ReadConfiguration()
	models.DBInit()
	iris.Config.MaxRequestBodySize = models.Config.MaxUploadFileSize //1G
	iris.Config.IsDevelopment = true
	iris.Get("/", func(ctx *iris.Context) {
		ctx.MustRender("index.html", nil)
	})
	iris.Static("/public", "./public/", 1)
	iris.StaticFS("/upload", models.Config.UploadDir, 1)
	iris.Get("/msg", func(ctx *iris.Context) {
		pgStr := ctx.FormValueString("pagesize")
		nthStr := ctx.FormValueString("nth")
		filter := ctx.FormValueString("filter")

		if pgStr == "" && nthStr == "" {
			resp := EmptySuccessResponse
			count, err := models.GetMsgCount(filter)
			if err != nil {
				ctx.WriteString(err.Error())
				return
			}
			resp.Data = count
			ctx.JSON(200, resp)
			return
		}
		pagesize, err := strconv.Atoi(pgStr)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}
		nth, err := strconv.Atoi(nthStr)
		if err != nil {
			ctx.WriteString(err.Error())
			return
		}

		msg, err := models.SelectPageMsg(filter, pagesize, nth)
		if err != nil {
			ctx.WriteString(err.Error())
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
			ctx.WriteString(err.Error())
			return
		}

		switch msg.Op {
		case "INSERT":
			if err := models.InsertMsg(msg.Data, "admin"); err != nil {
				ctx.WriteString(err.Error())
				return
			}
			ctx.JSON(200, EmptySuccessResponse)
		case "DELETE":
			id, err := strconv.ParseInt(msg.Data, 10, 0)
			if err != nil {
				ctx.WriteString(err.Error())
				return
			}
			if err := models.DeleteMsg(id); err != nil {
				ctx.WriteString(err.Error())
				return
			}
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
			fname := models.GetFilePrefix() + header.Filename
			path := filepath.Join(models.Config.UploadDir, fname)
			fasthttp.SaveMultipartFile(header, path)
			files[header.Filename] = "/upload/" + fname
		}
	})
	iris.Listen(models.Config.Listen + ":" + models.Config.Port)
}
