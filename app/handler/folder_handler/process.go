package folder_handler

import (
	"github.com/gin-gonic/gin"
	"haoyu.love/ImageServer/app"
	"haoyu.love/ImageServer/app/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func Process(c *gin.Context) {
	path := c.Param("path")
	fullPath := filepath.Join(*app.Root, path)
	fullPath, _ = filepath.Abs(fullPath)

	if ok, _ := util.IsSubFolder(*app.Root, fullPath); !ok {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		c.String(http.StatusNotFound, "Not found")
	} else if fileInfo.IsDir() {
		processSingleFolder(c)
	} else {
		processFile(c)
	}
}

func processFile(c *gin.Context) {
	path := filepath.Join(*app.Root, c.Param("path"))
	c.File(path)
}

func processSingleFolder(c *gin.Context) {
	path := filepath.Join(*app.Root, c.Param("path"))

	folders, files, err := GetFolderContent(path)
	pageNumStr := c.DefaultQuery("p", "1")
	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}
	folders, files, pageNum, pageNumMax, pagePrev, pageNext := util.Pagination(*app.PageSize, pageNum, folders, files)

	folders = util.RemoveLeft(*app.Root, folders)
	files = util.RemoveLeft(*app.Root, files)

	navigation := gin.H{
		"path": c.Param("path"),
		"prev": "",
		"next": "",
	}

	pagination := gin.H{
		"num":  pageNum,
		"max":  pageNumMax,
		"prev": pagePrev,
		"next": pageNext,
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"folders":    folders,
		"files":      files,
		"pagination": pagination,
		"navigation": navigation,
	})
}
