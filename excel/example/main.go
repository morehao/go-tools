package main

import (
	"fmt"

	"github.com/morehao/golib/excel"
	"github.com/xuri/excelize/v2"
)

func main() {
	read()
	write()
}

type DataItem struct {
	SerialNumber int64  `ex:"head:序号" validate:"min=10,max=100"`
	UserName     string `ex:"head:姓名"`
	Age          int64  `ex:"head:年龄"`
}

func read() {
	f, openErr := excelize.OpenFile("test.xlsx")
	if openErr != nil {
		fmt.Println("open file error: ", openErr)
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
	}
	if len(validateErrMap) > 0 {
		fmt.Println("validate error: ", validateErrMap)
	}
	for _, item := range dataList {
		fmt.Println(item)
	}
}

func write() {
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
