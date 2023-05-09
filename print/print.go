package myPrint

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

var printLabels []string

// 定义BarTender应用程序变量和标签文档变量
var btApp *ole.IDispatch

// var btFormatDocument *ole.IDispatch
var btFormatDocument *ole.VARIANT

func PrintTool(filePath string, list []map[string]interface{}) []string {

	//异常捕获
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("打印机报错", err)
			fmt.Println("标签文档关闭")
			// 关闭标签文档
			if btFormatDocument != nil {
				_, err := btFormatDocument.ToIDispatch().CallMethod("Close", 0)
				if err != nil {
					panic(err)
				}
				btFormatDocument.ToIDispatch().Release()
			}

			fmt.Println("退出BarTender应用程序")
			// 退出BarTender应用程序
			if btApp != nil {
				_, err := btApp.CallMethod("Quit", 0)
				if err != nil {
					panic(err)
				}
				btApp.Release()
			}
		}
	}()

	// 初始化COM库
	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	if err != nil {
		panic(err)
	}

	defer ole.CoUninitialize()

	// 创建BarTender应用程序实例
	unknown, err := oleutil.CreateObject("BarTender.Application")
	if err != nil {
		panic(err)
	}
	defer unknown.Release()

	// 获取IDispatch接口实例，用于调用COM组件方法
	btApp, err = unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		panic(err)
	}

	property, err := btApp.GetProperty("Formats")
	if err != nil {
		panic(err)
	}

	// 打开标签文档
	btFormatDocument, err = oleutil.CallMethod(property.ToIDispatch(), "Open", filePath, false, "")
	if err != nil {
		panic(err)
	}
	//启动打印驱动
	method, err := oleutil.CallMethod(btFormatDocument.ToIDispatch(), "PrintSetup")
	//数据填充&执行打印
	for _, mapKey := range list {

		for s := range mapKey {
			_, err := oleutil.CallMethod(btFormatDocument.ToIDispatch(), "SetNamedSubStringValue", s, mapKey[s])
			if err != nil {
				panic(err)
			}
		}

		oleutil.CallMethod(method.ToIDispatch(), "IdenticalCopiesOfLabel", 1)

		// 执行打印操作
		_, err = oleutil.CallMethod(btFormatDocument.ToIDispatch(), "PrintOut", false, false)
		if err != nil {
			panic(err)
		}
		val := mapKey["labelCode"]
		str, _ := val.(string)
		printLabels = append(printLabels, str)
	}

	//程序打印
	defer func() {
		fmt.Println("标签文档关闭")
		// 关闭标签文档
		if btFormatDocument != nil {
			_, err := btFormatDocument.ToIDispatch().CallMethod("Close", 0)
			if err != nil {
				panic(err)
			}
			btFormatDocument.ToIDispatch().Release()
		}

		fmt.Println("退出BarTender应用程序")
		// 退出BarTender应用程序
		if btApp != nil {
			_, err := btApp.CallMethod("Quit", 0)
			if err != nil {
				panic(err)
			}
			btApp.Release()
		}
	}()

	return printLabels

}

type Win32Printer struct {
	Name        string //名称
	WorkOffline bool   //是否在线
	Default     bool
	Status      bool
}

// GetPrintStatus 获取打印机状态
func GetPrintStatus(printerName string) bool {
	onLine := false

	var dst []Win32Printer
	q := "SELECT * FROM Win32_Printer"
	if err := wmi.Query(q, &dst); err != nil {
		panic(err)
	}

	for _, printer := range dst {
		fmt.Printf("打印机名称：%s\n", printer.Name)
		fmt.Printf("是否在线：%v\n", printer.WorkOffline)
		fmt.Println()

		if printer.Name == printerName && printer.WorkOffline == false {
			onLine = true
		}
	}

	return onLine

}

func GetPrinter() []map[string]interface{} {

	var dst []Win32Printer
	q := "SELECT * FROM Win32_Printer"
	if err := wmi.Query(q, &dst); err != nil {
		panic(err)
	}

	array := make([]map[string]interface{}, 0)

	for _, printer := range dst {

		printerMap := make(map[string]interface{})

		printerMap["name"] = printer.Name
		printerMap["default"] = printer.Default
		printerMap["status"] = printer.Status
		array = append(array, printerMap)
	}

	return array

}
