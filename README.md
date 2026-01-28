# go-docx2pdf

一个高性能的原生 Golang 库，用于将 DOCX 转换为 PDF，无需依赖 LibreOffice 等外部工具。

## 说明

本项目基于 [ryugenxd/docx2pdf](https://github.com/ryugenxd/docx2pdf) 二次修改，主要修复了 Go 版本兼容性问题，现已支持 **Go 1.20** 版本。

## 特性

- 纯 Go 实现，无外部依赖
- 高性能转换
- 支持 Go 1.20+

## 安装

```bash
go get github.com/26552594a/go-docx2pdf
```

## 使用示例

```go
package main

import (
	"fmt"
	"log"

	"github.com/26552594a/go-docx2pdf"
)

func main() {
	// 将 DOCX 文件转换为 PDF
	err := docx2pdf.ConvertFile("./input.docx", "./output.pdf")
	if err != nil {
		log.Fatal("转换失败:", err)
	}
	fmt.Println("转换成功！")
}
```

## 致谢

感谢原作者 [ryugenxd](https://github.com/ryugenxd) 的开源。
