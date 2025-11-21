package pdf

import (
	"fmt"
	"project/internal/models"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func CreateBeautifulStructPDF(responses []models.LinksResponce) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(0, 0, 128)
	pdf.Cell(0, 10, "Links Document")
	pdf.Ln(15)

	for i, response := range responses {
		pdf.SetFont("Arial", "B", 14)
		pdf.SetTextColor(0, 100, 0)
		pdf.Cell(0, 8, fmt.Sprintf("Links Response #%d:", i+1))
		pdf.Ln(8)

		pdf.SetFont("Courier", "", 10)
		pdf.SetTextColor(0, 0, 0)

		example := generateStructString(response, i+1)
		pdf.MultiCell(0, 5, example, "", "", false)

		if i < len(responses)-1 {
			pdf.Ln(10)
		}
	}

	return pdf
}

func generateStructString(response models.LinksResponce, index int) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("LinksResponse%d{\n", index))
	builder.WriteString("    Links: map[string]string{\n")

	for key, value := range response.Links {
		builder.WriteString(fmt.Sprintf("        \"%s\": \"%s\",\n", key, value))
	}

	builder.WriteString("    },\n")
	builder.WriteString(fmt.Sprintf("    LinksNum: %d,\n", response.LinksNum))
	builder.WriteString("}")

	return builder.String()
}
