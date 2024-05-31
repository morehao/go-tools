package autoCode

import (
	"os"
	"text/template"
)

func Tmpl() {
	// 模板定义
	tepl := "My name is {{.}}"
	// 解析模板
	tmpl, err := template.New("test").Parse(tepl)
	if err != nil {
		panic(err)
	}
	// 渲染模板
	if err := tmpl.Execute(os.Stdout, "Jack"); err != nil {
		panic(err)
	}

}
