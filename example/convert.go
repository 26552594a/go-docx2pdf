package main

import (
	"log"

	"github.com/ryugenxd/docx2pdf"
)

func main() {
	err := docx2pdf.ConvertFile("./TestDocument.docx", "./out.pdf")
	if err != nil {
		log.Fatal(err)
	}
}
