package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goacs/lib"
	"net/http"
	"os"
	"path/filepath"
)

func ListFiles(ctx *gin.Context) {

}

func UploadFile(ctx *gin.Context) {
	env := lib.Env{}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileDir := env.Get("FILESTORE_PATH", "./storage")
	absPath, err := filepath.Abs(fileDir)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("upload file err: %s", err.Error())})
		return
	}

	filePath := filepath.Join(absPath, filepath.Base(fileDir+"/"+file.Filename))

	if fileExists(filePath) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File %s exists", file.Filename)})
		return
	}

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("upload file err: %s", err.Error())})
		return
	}
}

func DownloadFile(ctx *gin.Context) {

}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
