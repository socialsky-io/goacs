package controllers

import (
	"github.com/gin-gonic/gin"
	"goacs/http/response"
	"goacs/models/templates"
	"goacs/repository"
	"goacs/repository/mysql"
	"strconv"
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

func GetTemplate(ctx *gin.Context) {
	templateId, err := strconv.Atoi(ctx.Param("templateid"))

	if err != nil {
		response.ResponseError(ctx, 404, "", "")
		return
	}

	templatesrepository := mysql.NewTemplateRepository(repository.GetConnection())
	template, err := templatesrepository.Find(int64(templateId))

	if err != nil {
		response.ResponseError(ctx, 404, "", "")
		return
	}

	response.ResponseData(ctx, template)
}

func GetTemplateParameters(ctx *gin.Context) {
	templateId, err := strconv.Atoi(ctx.Param("templateid"))

	paginatorRequest := repository.PaginatorRequestFromContext(ctx)
	templatesrepository := mysql.NewTemplateRepository(repository.GetConnection())
	template, err := templatesrepository.Find(int64(templateId))

	if err == nil {
		parameters, total := templatesrepository.ListTemplateParameters(template, paginatorRequest)
		responseData := repository.NewPaginatorResponse(paginatorRequest, total, parameters)
		response.ResponsePaginatior(ctx, responseData)
		return
	}

}
