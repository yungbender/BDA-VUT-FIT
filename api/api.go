package api

import (
	"bda/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code uint        `json:"code"`
	Data interface{} `json:"data"`
}

type StatusResponse struct {
	OnlineNodes  int `json:"online"`
	OfflineNodes int `json:"offline"`
	UnknownNodes int `json:"unknown"`
}

type StatusSelect struct {
	Active int
	Cnt    int
}

func StatusHandler(ctx *gin.Context) {
	db, err := models.GetDb()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Data: struct{}{}})
		return
	}

	var sel []StatusSelect
	db.Table("nodes").Select("active, COUNT(*) AS cnt").Group("active").Find(&sel)

	var res StatusResponse
	for i := range sel {
		r := sel[i]

		switch r.Active {
		case 0:
			res.OfflineNodes = r.Cnt
		case 1:
			res.OnlineNodes = r.Cnt
		case 2:
			res.UnknownNodes = r.Cnt
		}
	}
	ctx.JSON(http.StatusOK, Response{Code: http.StatusOK, Data: res})
	return
}

func Start() {
	router := gin.Default()

	router.GET("/v1/status", StatusHandler)
	router.Run(":13337")
}
