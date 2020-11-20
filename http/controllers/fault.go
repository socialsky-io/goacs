package controllers

import (
	"github.com/gin-gonic/gin"
	"goacs/http/response"
	"goacs/repository/mysql"
)

func GetTodayFaults(ctx *gin.Context) {
	faultRepository := mysql.NewFaultRepository()
	faults := faultRepository.GetLastDay(100)
	response.ResponseData(ctx, faults)
}
