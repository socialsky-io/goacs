package controllers

import (
	"github.com/gin-gonic/gin"
	"goacs/http/response"
	"goacs/repository"
	"goacs/repository/mysql"
)

func GetTemplatesList(ctx *gin.Context) {
	paginatorRequest := repository.PaginatorRequestFromContext(ctx)
	templatesrepository := mysql.NewTemplateRepository(repository.GetConnection())
	templates, total := templatesrepository.List(paginatorRequest)
	responseData := repository.NewPaginatorResponse(paginatorRequest, total, templates)
	response.ResponsePaginatior(ctx, responseData)
}
