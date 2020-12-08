package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goacs/http/response"
	"goacs/lib"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileInfoResponse struct {
	Size     int64     `json:"size"`
	Filename string    `json:"filename"`
	IsDir    bool      `json:"is_dir"`
	ModTime  time.Time `json:"mod_time"`
}

func ListFiles(ctx *gin.Context) {
	env := lib.Env{}

	fileDir := env.Get("FILESTORE_PATH", "./storage")
	absPath, _ := filepath.Abs(fileDir)
	files, err := ioutil.ReadDir(absPath)

	if err != nil {
		response.ResponseError(ctx, http.StatusInternalServerError, "File list error", err)
		return
	}

	var fileResponse []FileInfoResponse

	for _, file := range files {
		fileResponse = append(fileResponse, FileInfoResponse{
			Size:     file.Size(),
			Filename: file.Name(),
			IsDir:    file.IsDir(),
			ModTime:  file.ModTime(),
		})
	}
	response.ResponseData(ctx, fileResponse)
}

func UploadFile(ctx *gin.Context) {
	env := lib.Env{}

	file, err := ctx.FormFile("file")
	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, "", err)
		return
	}

	fileDir := env.Get("FILESTORE_PATH", "./storage")
	absPath, err := filepath.Abs(fileDir)

	if err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()), err)
		return
	}

	filePath := filepath.Join(absPath, filepath.Base(fileDir+"/"+file.Filename))

	if fileExists(filePath) {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Sprintf("File %s exists", file.Filename), err)
		return
	}

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		response.ResponseError(ctx, http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()), err)
		return
	}
}

func DownloadFile(ctx *gin.Context) {
	env := lib.Env{}
	fileDir := env.Get("FILESTORE_PATH", "./storage")
	absPath, _ := filepath.Abs(fileDir)
	fileName := ctx.Param("filename")
	//Currently dangerous method :/
	//TODO: add filepath security check
	filePath := filepath.Join(absPath, filepath.Base(fileDir+"/"+fileName))

	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName)) //fmt.Sprintf("attachment; filename=%s", filename) Downloaded file renamed
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(filePath)

}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
