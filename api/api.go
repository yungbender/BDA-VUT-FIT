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

type UserAgentsResponse []string

type UserAgentsSelect struct {
	UserAgent string
}

func UserAgentsHandler(ctx *gin.Context) {
	db, err := models.GetDb()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Data: struct{}{}})
		return
	}

	var uas []UserAgentsSelect
	var res UserAgentsResponse
	db.Table("nodes").Select("DISTINCT user_agent").Where("user_agent != ?", "").Find(&uas)

	for i := range uas {
		ua := uas[i]

		res = append(res, ua.UserAgent)
	}

	ctx.JSON(http.StatusOK, Response{Code: http.StatusOK, Data: res})
	return
}

func Start() {
	router := gin.Default()

	router.GET("/v1/status", StatusHandler)
	router.GET("/v1/useragents", UserAgentsHandler)
	router.Run(":13337")
}
