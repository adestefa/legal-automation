package services

import (
	"fmt"
	"strings"
)

// LegalDocumentFormatter handles legal document formatting and styling
type LegalDocumentFormatter struct {
	Style DocumentStyle
}

// DocumentStyle defines formatting preferences
type DocumentStyle struct {
	LineSpacing     float64
	FontSize        int
	FontFamily      string
	Margins         Margins
	PageSize        string
	HeaderStyle     string
	ParagraphStyle  string
}

// Margins defines document margins
type Margins struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
}

// NewLegalDocumentFormatter creates a new document formatter
func NewLegalDocumentFormatter() *LegalDocumentFormatter {
	return &LegalDocumentFormatter{
		Style: DocumentStyle{
			LineSpacing:    1.5,
			FontSize:       12,
			FontFamily:     "Times New Roman",
			PageSize:       "Letter",
			HeaderStyle:    "center",
			ParagraphStyle: "justified",
			Margins: Margins{
				Top:    1.0,
				Bottom: 1.0,
				Left:   1.0,
				Right:  1.0,
			},
		},
	}
}

// FormatAsHTML formats the document content as HTML for web display
func (ldf *LegalDocumentFormatter) FormatAsHTML(content string) string {
	var html strings.Builder
	
	html.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Legal Document</title>
    <style>
        body {
            font-family: 'Times New Roman', serif;
            font-size: 12pt;
            line-height: 1.5;
            margin: 1in;
            color: #000;
            background: #fff;
        }
        .document-header {
            text-align: center;
            margin-bottom: 2em;
            font-weight: bold;
        }
        .section-header {
            text-align: center;
            margin: 1.5em 0 1em 0;
            font-weight: bold;
            text-decoration: underline;
        }
        .paragraph {
            margin-bottom: 1em;
            text-align: justify;
            text-indent: 0.5in;
        }
        .parties-section {
            margin: 2em 0;
        }
        .signature-block {
            margin-top: 3em;
            text-align: left;
        }
        .page-break {
            page-break-before: always;
        }
        @media print {
            body { margin: 1in; }
            .page-break { page-break-before: always; }
        }
    </style>
</head>
<body>
    <div class="legal-document">`)
	
	// Process content sections
	sections := strings.Split(content, "\n\n")
	for i, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}
		
		// Determine section type and format accordingly
		if ldf.isHeaderSection(section) {
			html.WriteString(fmt.Sprintf(`        <div class="document-header">%s</div>`, ldf.formatTextForHTML(section)))
		} else if ldf.isSectionHeader(section) {
			html.WriteString(fmt.Sprintf(`        <div class="section-header">%s</div>`, ldf.formatTextForHTML(section)))
		} else if ldf.isSignatureBlock(section) {
			html.WriteString(fmt.Sprintf(`        <div class="signature-block">%s</div>`, ldf.formatTextForHTML(section)))
		} else {
			// Regular paragraph
			html.WriteString(fmt.Sprintf(`        <div class="paragraph">%s</div>`, ldf.formatTextForHTML(section)))
		}
		
		// Add page breaks where appropriate
		if ldf.shouldAddPageBreak(section, i, len(sections)) {
			html.WriteString(`        <div class="page-break"></div>`)
		}
	}
	
	html.WriteString(`    </div>
</body>
</html>`)
	
	return html.String()
}

// FormatAsPDF formats the document for PDF generation (placeholder)
func (ldf *LegalDocumentFormatter) FormatAsPDF(content string) ([]byte, error) {
	// This would integrate with a PDF generation library
	// For now, return HTML that can be converted to PDF
	htmlContent := ldf.FormatAsHTML(content)
	return []byte(htmlContent), nil
}

// FormatAsPlainText formats the document as plain text
func (ldf *LegalDocumentFormatter) FormatAsPlainText(content string) string {
	// Clean up the content for plain text display
	formatted := strings.ReplaceAll(content, "\n\n\n", "\n\n")
	formatted = strings.TrimSpace(formatted)
	
	// Add proper spacing for legal document structure
	sections := strings.Split(formatted, "\n\n")
	var result strings.Builder
	
	for i, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}
		
		// Add appropriate spacing
		if i > 0 {
			result.WriteString("\n\n")
		}
		
		// Format section based on type
		if ldf.isHeaderSection(section) {
			// Center header text (approximation for plain text)
			lines := strings.Split(section, "\n")
			for _, line := range lines {
				result.WriteString(ldf.centerText(line, 80))
				result.WriteString("\n")
			}
		} else {
			result.WriteString(section)
		}
	}
	
	return result.String()
}

// Helper methods

func (ldf *LegalDocumentFormatter) isHeaderSection(text string) bool {
	upperText := strings.ToUpper(text)
	return strings.Contains(upperText, "UNITED STATES DISTRICT COURT") ||
		   strings.Contains(upperText, "COMPLAINT FOR") ||
		   strings.Contains(upperText, "CASE NO.")
}

func (ldf *LegalDocumentFormatter) isSectionHeader(text string) bool {
	upperText := strings.ToUpper(strings.TrimSpace(text))
	sectionHeaders := []string{
		"PARTIES",
		"FACTUAL ALLEGATIONS", 
		"CAUSES OF ACTION",
		"DAMAGES",
		"PRAYER FOR RELIEF",
		"COUNT I",
		"COUNT II",
		"COUNT III",
		"COUNT IV",
		"COUNT V",
	}
	
	for _, header := range sectionHeaders {
		if upperText == header || strings.HasPrefix(upperText, header) {
			return true
		}
	}
	
	return false
}

func (ldf *LegalDocumentFormatter) isSignatureBlock(text string) bool {
	return strings.Contains(text, "Respectfully submitted") ||
		   strings.Contains(text, "_________________________") ||
		   strings.Contains(text, "Attorney for Plaintiff")
}

func (ldf *LegalDocumentFormatter) shouldAddPageBreak(section string, index int, total int) bool {
	// Add page breaks before major sections
	upperText := strings.ToUpper(strings.TrimSpace(section))
	
	// Break before causes of action if document is getting long
	if strings.Contains(upperText, "CAUSES OF ACTION") && index > 2 {
		return true
	}
	
	// Break before prayer for relief
	if strings.Contains(upperText, "PRAYER FOR RELIEF") {
		return true
	}
	
	return false
}

func (ldf *LegalDocumentFormatter) formatTextForHTML(text string) string {
	// Escape HTML characters and preserve line breaks
	escaped := strings.ReplaceAll(text, "&", "&amp;")
	escaped = strings.ReplaceAll(escaped, "<", "&lt;")
	escaped = strings.ReplaceAll(escaped, ">", "&gt;")
	escaped = strings.ReplaceAll(escaped, "\"", "&quot;")
	escaped = strings.ReplaceAll(escaped, "'", "&#39;")
	
	// Convert line breaks to HTML
	escaped = strings.ReplaceAll(escaped, "\n", "<br>")
	
	return escaped
}

func (ldf *LegalDocumentFormatter) centerText(text string, width int) string {
	text = strings.TrimSpace(text)
	if len(text) >= width {
		return text
	}
	
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

// AddHighlighting adds highlighting to specific terms in the document
func (ldf *LegalDocumentFormatter) AddHighlighting(content string, terms []string) string {
	highlighted := content
	
	for _, term := range terms {
		if term == "" {
			continue
		}
		
		// Create highlight span
		highlightStart := fmt.Sprintf(`<span class="highlight" style="background-color: yellow; font-weight: bold;">%s</span>`, term)
		
		// Replace all occurrences (case insensitive)
		highlighted = strings.ReplaceAll(highlighted, term, highlightStart)
		highlighted = strings.ReplaceAll(highlighted, strings.ToUpper(term), highlightStart)
		highlighted = strings.ReplaceAll(highlighted, strings.ToLower(term), highlightStart)
		highlighted = strings.ReplaceAll(highlighted, strings.Title(term), highlightStart)
	}
	
	return highlighted
}

// GenerateDocumentMetrics calculates various document metrics
func (ldf *LegalDocumentFormatter) GenerateDocumentMetrics(content string) DocumentMetrics {
	words := strings.Fields(content)
	paragraphs := strings.Split(content, "\n\n")
	pages := ldf.estimatePageCount(content)
	
	return DocumentMetrics{
		WordCount:      len(words),
		ParagraphCount: len(paragraphs),
		PageCount:      pages,
		CharacterCount: len(content),
		ReadingTime:    ldf.estimateReadingTime(len(words)),
	}
}

// DocumentMetrics contains document statistics
type DocumentMetrics struct {
	WordCount      int     `json:"wordCount"`
	ParagraphCount int     `json:"paragraphCount"`
	PageCount      int     `json:"pageCount"`
	CharacterCount int     `json:"characterCount"`
	ReadingTime    float64 `json:"readingTimeMinutes"`
}

func (ldf *LegalDocumentFormatter) estimatePageCount(content string) int {
	// Estimate based on average words per page (250-300 for legal documents)
	words := len(strings.Fields(content))
	return (words / 275) + 1
}

func (ldf *LegalDocumentFormatter) estimateReadingTime(wordCount int) float64 {
	// Average reading speed for legal documents is slower (150-200 WPM)
	return float64(wordCount) / 175.0
}