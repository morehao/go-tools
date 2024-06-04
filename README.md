# go-tools
`go-tools`是一个golang工具库，包含了一些个人在项目开发过程中总结的一些常用的工具函数和组件。

## autoCode

### 简介
`autoCode` 是一个简单的代码生成工具，通过读取数据库表结构，支持生成基础的CRUD代码，router、controller、service、dto、model、errorCode等代码。
### 特性
- 支持MySQL数据库
- 支持模板自定义和模板参数自定义
- 支持生成基础的CRUD代码

### 安装
```bash
go get github.com/morehao/go-tools
```
### 使用
使用示例参照[autoCode单测](./autoCode/auto_code_test.go)

## excel

### 简介
`excel` 是基于 `excelize` 的简单封装，支持通过结构体便捷地读写 Excel 文件。

无论是读取 Excel 还是写入 Excel，都需要定义一个结构体，结构体的字段通过 tag（即 `ex`）来指定 Excel 的相关信息。

### 特性
- 通过结构体标签定义Excel列映射关系
- 支持读取和写入Excel文件
- 支持基于validator的数据验证

### 安装

```bash
go get github.com/morehao/go-tools
```
### 使用
### 读取Excel
读取Excel的简单示例：
```go 
package main

import (
	"fmt"
	"github.com/morehao/go-tools"
	"github.com/xuri/excelize/v2"
)

type DataItem struct {
	SerialNumber int64  `ex:"head:序号" validate:"min=10,max=100"`
	UserName     string `ex:"head:姓名"`
	Age          int64  `ex:"head:年龄"`
}

func main() {
	f, openErr := excelize.OpenFile("test.xlsx")
	if openErr != nil {
		fmt.Println("open file error: ", openErr)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	reader := excel.NewReader(f, &excel.ReaderOption{
		SheetNumber:  0,
		HeadRow:      0,
		DataStartRow: 1,
	})
	var dataList []DataItem
	validateErrMap, readerErr := reader.Read(&dataList)
	if readerErr != nil {
		fmt.Println("read error: ", readerErr)
		return
	}
	if len(validateErrMap) > 0 {
		fmt.Println("validate error: ", validateErrMap)
		return
	}
	for _, item := range dataList {
		fmt.Println(item)
	}
}
```
### 写入Excel
生成Excel的简单示例：
```go
package main

import (
	"fmt"
	"github.com/morehao/go-tools"
)

type DataItem struct {
	SerialNumber int64  `ex:"head:序号" validate:"min=10,max=100"`
	UserName     string `ex:"head:姓名"`
	Age          int64  `ex:"head:年龄"`
}

func main() {
	var dataList []DataItem
	dataList = append(dataList, DataItem{
		SerialNumber: 1,
		UserName:     "张三",
		Age:          18,
	})
	excelWriter := excel.NewWrite(&excel.WriteOption{
		SheetName: "Sheet1",
		HeadRow:   0,
	})
	if err := excelWriter.SaveAs(dataList, "write.xlsx"); err != nil {
		fmt.Println("write error: ", err)
	}
}
```
