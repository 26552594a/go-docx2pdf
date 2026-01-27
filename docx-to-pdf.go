package docx2pdf

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type Text struct {
	Text string `xml:",chardata"`
}

type RunProperties struct {
	Bold      bool   `xml:"b"`
	Italic    bool   `xml:"i"`
	FontSize  string `xml:"sz"`
	FontColor string `xml:"color"`
}

type Run struct {
	Properties RunProperties `xml:"rPr"`
	Texts      []Text        `xml:"t"`
}

type Paragraph struct {
	Alignment string `xml:"pPr>jc"` // Align: left, right, center
	Runs      []Run  `xml:"r"`
}

type TableCell struct {
	Text string `xml:"p>r>t"`
}

type TableRow struct {
	Cells []TableCell `xml:"tc"`
}

type Table struct {
	Rows []TableRow `xml:"tr"`
}

type Drawing struct {
	Image ImageRef `xml:"inline>graphic>graphicData>pic:pic>blipFill>blip"`
}

type ImageRef struct {
	ID string `xml:"embed,attr"`
}

type Document struct {
	Paragraphs []Paragraph `xml:"body>p"`
	Tables     []Table     `xml:"body>tbl"`
	Drawings   []Drawing   `xml:"body>drawing"`
}

func extractTextFromDocx(docxPath string) (string, error) {
	reader, err := zip.OpenReader(docxPath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var documentXML string
	for _, file := range reader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			bytes, err := ioutil.ReadAll(rc)
			if err != nil {
				return "", err
			}
			documentXML = string(bytes)
			break
		}
	}

	return documentXML, nil
}

func setFontFromRun(pdf *gofpdf.Fpdf, run Run) {
	fontStyle := ""
	if run.Properties.Bold {
		fontStyle += "B"
	}
	if run.Properties.Italic {
		fontStyle += "I"
	}
	fontSize := 12.0 // default font size
	if run.Properties.FontSize != "" {
		// convert font size from half-points to points
		fontSizeValue, err := strconv.ParseFloat(run.Properties.FontSize, 64)
		if err == nil {
			fontSize = fontSizeValue / 2
		}
	}
	pdf.SetFont("Arial", fontStyle, fontSize)
	if run.Properties.FontColor != "" {
		pdf.SetTextColor(parseHexColor(run.Properties.FontColor))
	}
}

func parseHexColor(s string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

func setParagraphAlignment(_ *gofpdf.Fpdf, alignment string) string {
	switch alignment {
	case "center":
		return "C"
	case "right":
		return "R"
	default:
		return "L"
	}
}

func processParagraph(pdf *gofpdf.Fpdf, para Paragraph) {
	align := setParagraphAlignment(pdf, para.Alignment)
	pdf.SetFont("Arial", "", 12)

	for _, run := range para.Runs {
		setFontFromRun(pdf, run)
		for _, text := range run.Texts {
			pdf.CellFormat(0, 6, text.Text, "", 1, align, false, 0, "")
		}
	}

	pdf.Ln(4) // Spasi antar paragraf
}

func processTable(pdf *gofpdf.Fpdf, table Table) {
	pdf.SetFont("Arial", "", 10)

	for _, row := range table.Rows {
		for _, cell := range row.Cells {
			pdf.CellFormat(40, 10, cell.Text, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1) // Pindah ke baris berikutnya
	}
}

func addImageToPDF(pdf *gofpdf.Fpdf, imgPath string, x, y, width, height float64) {
	pdf.Image(imgPath, x, y, width, height, false, "", 0, "")
}

func extractImagesFromDocx(_ string, reader *zip.ReadCloser) (map[string]string, error) {
	images := make(map[string]string)
	tempDir, err := ioutil.TempDir("", "docx_images")
	if err != nil {
		return nil, err
	}

	for _, file := range reader.File {
		if strings.HasPrefix(file.Name, "word/media/") {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			imgBytes, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			imgPath := filepath.Join(tempDir, filepath.Base(file.Name))
			err = ioutil.WriteFile(imgPath, imgBytes, 0644)
			if err != nil {
				return nil, err
			}

			images[file.Name] = imgPath
		}
	}
	return images, nil
}

func createPDF(text string, outputPath string, images map[string]string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	var doc Document
	err := xml.Unmarshal([]byte(text), &doc)
	if err != nil {
		return err
	}

	for _, para := range doc.Paragraphs {
		processParagraph(pdf, para)
	}

	for _, table := range doc.Tables {
		processTable(pdf, table)
	}

	for _, drawing := range doc.Drawings {
		imgPath, exists := images["word/media/"+drawing.Image.ID]
		if exists {
			addImageToPDF(pdf, imgPath, 10, 10, 50, 50) // Example coordinates and size
		}
	}

	err = pdf.OutputFileAndClose(outputPath)
	return err
}

func ConvertFile(inputFile string, outputFile string) error {
	reader, err := zip.OpenReader(inputFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	text, err := extractTextFromDocx(inputFile)
	if err != nil {
		return err
	}

	images, err := extractImagesFromDocx(inputFile, reader)
	if err != nil {
		return err
	}

	err = createPDF(text, outputFile, images)
	if err != nil {
		return err
	}

	return nil
}
