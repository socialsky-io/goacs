package controllers

import (
	"github.com/gin-gonic/gin"
	"goacs/http/response"
	"goacs/models/templates"
	"goacs/repository"
	"goacs/repository/mysql"
)

type TemplateListResponse struct {
	templates.Template
	ParameterCount int64 `json:"parameter_count"`
}

func GetTemplatesList(ctx *gin.Context) {
	paginatorRequest := repository.PaginatorRequestFromContext(ctx)
	templatesrepository := mysql.NewTemplateRepository(repository.GetConnection())
	templatesList, total := templatesrepository.List(paginatorRequest)
	templatesList = templatesrepository.HydrateTemplatesParameters(templatesList)

	var templateResponse []TemplateListResponse

	for _, template := range templatesList {
		count := int64(len(template.Parameters))
		template.Parameters = []templates.TemplateParameter{}
		templateResponse = append(templateResponse, TemplateListResponse{
			Template:       template,
			ParameterCount: count,
		})
	}

	responseData := repository.NewPaginatorResponse(paginatorRequest, total, templateResponse)
	response.ResponsePaginatior(ctx, responseData)
}
