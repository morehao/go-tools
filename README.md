# go-tools简介
`go-tools`是一个golang工具组件库，包含了一些个人在项目开发过程中总结的一些常用的工具函数和组件。

组件列表：
- [codeGen](#codegen) 代码生成工具
- `conc` 简单的并发控制组件
- `conf` 配置文件读取组件
- `dbClient` 数据库组件
- [excel](#excel) 简单读写excel组件
- `gast` 语法树工具
- `gcontext` 上下文工具组件
- `gerror` 错误处理组件
- `glog` 日志组件
- `gutils` 一些常用的工具函数
- `jwtAuth` jwt鉴权组件

# 安装
```bash
go get github.com/morehao/go-tools
```

# 组件使用说明

## codeGen

### 简介
`codeGen` 是一个简单的代码生成工具，通过读取数据库表结构，支持生成基础的CRUD代码，router、controller、service、dto、model、errorCode等代码。
### 特性
- 支持MySQL数据库
- 支持模板自定义和模板参数自定义
- 支持基于模板生成代码
### 使用
使用示例参照[codeGen单测](codeGen/gen_test.go)

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
使用示例参照[excel使用说明](excel/README.md)
