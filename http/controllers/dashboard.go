package controllers

import (
	"github.com/gin-gonic/gin"
	"goacs/http/response"
	"goacs/models/fault"
	"goacs/repository"
	"goacs/repository/mysql"
	"log"
)

type dashboardResponse struct {
	DevicesCount int64         `json:"devices_count"`
	InformsCount int64         `json:"informs_count"`
	FaultsCount  int64         `json:"faults_count"`
	Faults       []fault.Fault `json:"faults"`
}

func GetDashboardData(ctx *gin.Context) {
	cpeRepository := mysql.NewCPERepository(repository.GetConnection())
	fRepository := mysql.NewFaultRepository()

	responseData := dashboardResponse{
		DevicesCount: cpeRepository.Count(),
		InformsCount: 0,
		FaultsCount:  fRepository.CountLastDay(),
		Faults:       fRepository.GetLastDay(100),
	}

	log.Println(responseData.Faults)

	response.ResponseData(ctx, responseData)

}
