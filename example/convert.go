package main

import (
	"fmt"
	"log"

	"github.com/26552594a/go-docx2pdf"
)

func main() {
	// 将 DOCX 文件转换为 PDF
	err := docx2pdf.ConvertFile("./TestDocument.docx", "./out.pdf")
	if err != nil {
		log.Fatal("转换失败:", err)
	}
	fmt.Println("转换成功！已生成 out.pdf")
}
