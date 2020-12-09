package controllers

import (
	"github.com/gin-gonic/gin"
	acshttp "goacs/acs/http"
	"goacs/acs/types"
	"goacs/http/request"
	"goacs/http/response"
	"goacs/models/cpe"
	"goacs/repository"
	"goacs/repository/mysql"
	"log"
	"net/http"
	"strconv"
)

type ParameterRequest struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type AddObjectRequest struct {
	Name string `json:"name" binding:"required"`
	Key  string `json:"key"`
}

type AssignTemplateToDeviceRequest struct {
	TemplateId int64 `json:"template_id" validate:"required"`
	Priority   int64 `json:"priority" validate:"required"`
}

func GetDevice(ctx *gin.Context) {
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, err := getCPEFromContext(ctx, cperepository)
	if err != nil {
		response.ResponseError(ctx, 404, "Not found", "")
		return
	}

	response.ResponseData(ctx, cpeModel)

}

func GetDeviceParameters(ctx *gin.Context) {
	paginatorRequest := repository.PaginatorRequestFromContext(ctx)
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, err := getCPEFromContext(ctx, cperepository)
	if err == nil {
		parameters, total := cperepository.ListCPEParameters(cpeModel, paginatorRequest)
		responseData := repository.NewPaginatorResponse(paginatorRequest, total, parameters)
		response.ResponsePaginatior(ctx, responseData)
	}
}

func GetDeviceTemplates(ctx *gin.Context) {
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, _ := getCPEFromContext(ctx, cperepository)
	templaterepository := mysql.NewTemplateRepository(repository.GetConnection())
	templates := templaterepository.GetTemplatesForCPE(cpeModel)
	response.ResponseData(ctx, templates)
}

func AssignTemplateToDevice(ctx *gin.Context) {
	var assignDeviceRequest AssignTemplateToDeviceRequest
	_ = ctx.BindJSON(&assignDeviceRequest)
	log.Println(assignDeviceRequest)
	validator := request.NewApiValidator(ctx, assignDeviceRequest)
	verr := validator.Validate()

	if verr != nil {
		response.ResponseValidationErrors(ctx, validator)
		return
	}

	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, _ := getCPEFromContext(ctx, cperepository)
	templaterepository := mysql.NewTemplateRepository(repository.GetConnection())
	err := templaterepository.AssignTemplateToDevice(cpeModel, assignDeviceRequest.TemplateId, assignDeviceRequest.Priority)

	if err != nil {
		response.Response500(ctx, "", err)
		return
	}

	response.ResponseData(ctx, "")
}

func UnassignTemplateFromDevice(ctx *gin.Context) {
	templateId, _ := strconv.Atoi(ctx.Param("template_id"))
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, _ := getCPEFromContext(ctx, cperepository)
	templaterepository := mysql.NewTemplateRepository(repository.GetConnection())

	err := templaterepository.UnassignTemplateFromDevice(cpeModel, int64(templateId))

	if err != nil {
		response.Response500(ctx, "", err)
		return
	}

	response.ResponseData(ctx, "")

}

func GetDevicesList(ctx *gin.Context) {
	paginatorRequest := repository.PaginatorRequestFromContext(ctx)
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpes, total := cperepository.List(paginatorRequest)
	responseData := repository.NewPaginatorResponse(paginatorRequest, total, cpes)
	response.ResponsePaginatior(ctx, responseData)
}

func UpdateParameter(ctx *gin.Context) {
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, err := getCPEFromContext(ctx, cperepository)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var parameterRequest ParameterRequest
	if err := ctx.ShouldBind(&parameterRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, err = cperepository.UpdateParameter(cpeModel, types.ParameterValueStruct{
		Name:  parameterRequest.Name,
		Value: parameterRequest.Value,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(204, "")
}

func CreateParameter(ctx *gin.Context) {

}

func getCPEFromContext(ctx *gin.Context, cpeRepository mysql.CPERepository) (*cpe.CPE, error) {
	cpeModel, err := cpeRepository.Find(ctx.Param("uuid"))

	if err != nil {
		ctx.AbortWithError(404, err)
		return nil, err
	}

	return cpeModel, nil
}

func GetParameterValues(ctx *gin.Context) {
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, err := getCPEFromContext(ctx, cperepository)

	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err.Error(), "")
		return
	}

	parameters, err := cperepository.GetCPEParameters(cpeModel)

	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err.Error(), "")
		return
	}

	acsRequest := acshttp.NewACSRequest(cpeModel)
	acsRequest.GetParameterValues(cpe.DetermineDeviceTreeRootPath(parameters))

}

func GetDeviceQueuedTasks(ctx *gin.Context) {
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, err := getCPEFromContext(ctx, cperepository)

	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, err.Error(), "")
		return
	}

	taskRepository := mysql.NewTasksRepository(repository.GetConnection())
	tasks := taskRepository.GetTasksForCPE(cpeModel.UUID)

	response.ResponseData(ctx, tasks)
}

func AddObject(ctx *gin.Context) {
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, err := getCPEFromContext(ctx, cperepository)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	addObjectRequest := AddObjectRequest{}
	err = ctx.ShouldBindJSON(&addObjectRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	acsRequest := acshttp.NewACSRequest(cpeModel)
	acsRequest.AddObject(addObjectRequest.Name)
}

func Kick(ctx *gin.Context) {
	cperepository := mysql.NewCPERepository(repository.GetConnection())
	cpeModel, err := getCPEFromContext(ctx, cperepository)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	acsRequest := acshttp.NewACSRequest(cpeModel)
	acsRequest.Kick()
}
