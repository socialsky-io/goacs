package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goacs/models/fault"
	"goacs/repository"
	"goacs/repository/mysql"
)

type dashboardResponse struct {
	DevicesCount    int64         `json:"devices_count"`
	InformsCount    int64         `json:"informs_count"`
	ExceptionsCount int64         `json:"exceptions_count"`
	Faults          []fault.Fault `json:"faults"`
}

func GetDashboardData(ctx *gin.Context) {
	cpeRepository := mysql.NewCPERepository(repository.GetConnection())
	fRepository := mysql.NewFaultRepository()

	response := dashboardResponse{
		DevicesCount:    cpeRepository.Count(),
		InformsCount:    0,
		ExceptionsCount: fRepository.CountLastDay(),
		Faults:          fRepository.GetLastDay(100),
	}

	json.NewEncoder(ctx.Writer).Encode(response)

}
