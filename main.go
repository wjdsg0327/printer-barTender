package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	_ "printer/myutil"
	myFileUtil "printer/myutil"
	_ "printer/print"
	myPrint "printer/print"
	"time"
)

type PrintParamVo struct {
	LicenceKey     string
	BtwTemplateUrl string
	Data           []map[string]interface{}
}

var tool []string

func main() {

	//创建一个服务
	ginServer := gin.Default()

	ginServer.POST("/backend/client/print/btwPrint", func(c *gin.Context) {

		//获取当前时间
		now := time.Now()
		fmt.Println("start=================", now.Format("2006-01-02 15:04:05"))

		var printParamVo []PrintParamVo

		err := c.ShouldBind(&printParamVo)
		if err != nil {
			fmt.Println("错误：", err)
		}

		for _, v := range printParamVo {
			//判断校验码
			if len(v.LicenceKey) == 0 {
				c.JSON(http.StatusBadRequest, "校验码不能为空!")
				return
			}
			//判断模板文件目录是否存在
			err := os.Mkdir("C:\\BTWPrintTemplate", os.ModePerm)
			if err != nil {
				fmt.Println("文件已存在！")
			}
			//获取文件名字
			name := myFileUtil.GetFileName(v.BtwTemplateUrl)
			fmt.Println("文件名字是", name)
			//模板下载
			btwFilePath := "C:\\BTWPrintTemplate\\" + name
			fileExists := myFileUtil.FileExists(btwFilePath)
			if !fileExists {
				err := myFileUtil.DownloadFile(v.BtwTemplateUrl, "C:\\BTWPrintTemplate\\")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": "文件下载失败",
						"err": err,
					})
					return
				}
			}
			tool = myPrint.PrintTool(btwFilePath, v.Data)

		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    "200",
			"message": "操作成功",
			"result":  tool,
		})

		//获取当前时间
		end := time.Now()
		fmt.Println("entTime=================", end.Format("2006-01-02 15:04:05"))

	})

	//获取指定打印机状态
	ginServer.GET("/backend/client/print/getPrintStatus", func(c *gin.Context) {

		printerName := c.Query("printerName")

		status := myPrint.GetPrintStatus(printerName)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    "200",
			"message": "操作成功",
			"result":  status,
		})

	})

	//打印机列表信息查询接口
	ginServer.GET("/backend/client/print/getPrinter", func(c *gin.Context) {

		printer := myPrint.GetPrinter()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    "200",
			"message": "操作成功",
			"result":  printer,
		})

	})

	//服务器端口
	ginServer.Run(":40100")
}
