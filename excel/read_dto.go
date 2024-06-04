package excel

type ValidationError struct {
	DataRowNumber int    // Excel 表中的行号，从0开始
	Head          string // Excel 表中的列名（即表头名）
	CellValue     string // Excel 表中的单元格值
	ExpectType    string // 期望的单元格类型
	ErrorMessage  string // 错误信息
}
