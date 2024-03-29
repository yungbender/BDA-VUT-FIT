package api

import (
	"bda/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) StatusHandler(ctx *gin.Context) {
	db := h.DB

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
}

type UserAgentsSelect struct {
	UserAgent  string  `json:"user_agent"`
	Percentage float64 `json:"percentage"`
	Count      int     `json:"count"`
}

type UserAgentsCountSelect struct {
	Total int
}

type UserAgentsResponse struct {
	UserAgents []UserAgentsSelect `json:"user_agents"`
	Total      int                `json:"total_count"`
}

func (h *Handler) UserAgentsHandler(ctx *gin.Context) {
	db := h.DB

	var uas []UserAgentsSelect
	var res UserAgentsResponse

	db.Table("nodes").
		Select("user_agent, ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) AS percentage, COUNT(*) AS count").
		Where("user_agent != ?", "").
		Group("user_agent").
		Find(&uas)

	for i := range uas {
		ua := uas[i]

		res.UserAgents = append(res.UserAgents, ua)
	}

	var t UserAgentsCountSelect
	db.Table("nodes").
		Select("COUNT(*) AS total").
		Where("user_agent != ?", "").
		Find(&t)

	res.Total = t.Total

	ctx.JSON(http.StatusOK, Response{Code: http.StatusOK, Data: res})
}

type LiveNodeObj struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type LiveNodesResponse []LiveNodeObj

func (h *Handler) LiveNodesHandler(ctx *gin.Context) {
	db := h.DB

	var nodes []models.Node
	db.Where("active = ?", 1).Find(&nodes)

	var res LiveNodesResponse
	for i := range nodes {
		node := nodes[i]

		res = append(res, LiveNodeObj{IP: node.Ip, Port: node.Port})
	}
	ctx.JSON(http.StatusOK, Response{Code: http.StatusOK, Data: res})
}

func Start() {
	router := gin.Default()

	db, err := models.GetDb()
	if err != nil {
		panic(err)
	}
	handler := Handler{DB: db}

	router.GET("/v1/status", handler.StatusHandler)
	router.GET("/v1/useragents", handler.UserAgentsHandler)
	router.GET("/v1/livenodes", handler.LiveNodesHandler)
	router.Run(":8080")
}
