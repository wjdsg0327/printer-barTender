# 启动教程

安装相关mod

```go
go mod tidy
```

打包

```go
go build main.go
```

# 接口说明：

## 1、打印接口：

### url

```apl
/backend/client/print/btwPrint
```

### 参数说明

```json
[
  {
    "licenceKey": "", //校验码
    "btwTemplateUrl": "", //模板地址（可以下载的地址，比如minio）
    "data": [] //需要打印的数据
  }
]
```

用于接收的结构体

```go
type PrintParamVo struct {
	LicenceKey     string
	BtwTemplateUrl string
	Data           []map[string]interface{}
}
```

## 2、获取指定打印机状态

### url

```apl
/backend/client/print/getPrintStatus
```

### 参数

```go
printerName //字符串类型，打印机名称
```

### 返回

```json
c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    "200",
			"message": "操作成功",
			"result":  status,
})
```

### 返回结构体说明

```go
type Win32Printer struct {
	Name        string //打印机名称
	WorkOffline bool   //是否在线
	Default     bool  //是否是默认打印机
	Status      bool  //状态
}
```

## 3、打印机列表信息查询接口

### url

```
/backend/client/print/getPrinter
```

