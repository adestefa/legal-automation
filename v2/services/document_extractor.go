package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nguyenthenguyen/docx"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// DocumentExtractor handles text extraction from various document formats
type DocumentExtractor struct {
	SupportedFormats []string
	MaxFileSize      int64
}

// ExtractedContent represents the result of document text extraction
type ExtractedContent struct {
	RawText    string                 `json:"rawText"`
	Metadata   map[string]interface{} `json:"metadata"`
	PageCount  int                    `json:"pageCount"`
	WordCount  int                    `json:"wordCount"`
	Error      error                  `json:"error,omitempty"`
	SourceFile string                 `json:"sourceFile"`
}

// ContentPattern defines a pattern for extracting specific information
type ContentPattern struct {
	Pattern     *regexp.Regexp
	FieldName   string
	Required    bool
	Processor   func(string) interface{}
}

// NewDocumentExtractor creates a new document extractor instance
func NewDocumentExtractor() *DocumentExtractor {
	return &DocumentExtractor{
		SupportedFormats: []string{".pdf", ".docx", ".txt"},
		MaxFileSize:      50 * 1024 * 1024, // 50MB max file size
	}
}

// ExtractText extracts text content from a document file
func (e *DocumentExtractor) ExtractText(filePath string) (*ExtractedContent, error) {
	log.Printf("[EXTRACTOR] Starting text extraction for: %s", filePath)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}
	
	// Check file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}
	
	if fileInfo.Size() > e.MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max: %d)", fileInfo.Size(), e.MaxFileSize)
	}
	
	// Get file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	
	// Extract based on file type
	var content *ExtractedContent
	switch ext {
	case ".pdf":
		content, err = e.extractFromPDF(filePath)
	case ".docx":
		content, err = e.extractFromDOCX(filePath)
	case ".txt":
		content, err = e.extractFromTXT(filePath)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
	
	if err != nil {
		log.Printf("[EXTRACTOR] Error extracting from %s: %v", filePath, err)
		return &ExtractedContent{
			Error:      err,
			SourceFile: filePath,
			Metadata:   map[string]interface{}{"error": err.Error()},
		}, err
	}
	
	// Add metadata
	content.SourceFile = filePath
	content.WordCount = len(strings.Fields(content.RawText))
	if content.Metadata == nil {
		content.Metadata = make(map[string]interface{})
	}
	content.Metadata["fileSize"] = fileInfo.Size()
	content.Metadata["lastModified"] = fileInfo.ModTime()
	content.Metadata["extractedAt"] = time.Now()
	
	log.Printf("[EXTRACTOR] Successfully extracted %d words from %s", content.WordCount, filePath)
	return content, nil
}

// extractFromPDF extracts text from PDF files
func (e *DocumentExtractor) extractFromPDF(filePath string) (*ExtractedContent, error) {
	log.Printf("[EXTRACTOR] Extracting text from PDF: %s", filePath)
	
	// Open PDF file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF file: %v", err)
	}
	defer file.Close()
	
	// Create PDF reader
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %v", err)
	}
	
	// Check if PDF is encrypted
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, fmt.Errorf("failed to check PDF encryption: %v", err)
	}
	
	if isEncrypted {
		// Try to decrypt with empty password
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil || !auth {
			return nil, fmt.Errorf("PDF is password protected")
		}
	}
	
	// Get number of pages
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get PDF page count: %v", err)
	}
	
	// Extract text from all pages
	var textContent strings.Builder
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			log.Printf("[EXTRACTOR] Warning: failed to get page %d: %v", i, err)
			continue
		}
		
		extractor, err := extractor.New(page)
		if err != nil {
			log.Printf("[EXTRACTOR] Warning: failed to create extractor for page %d: %v", i, err)
			continue
		}
		
		text, err := extractor.ExtractText()
		if err != nil {
			log.Printf("[EXTRACTOR] Warning: failed to extract text from page %d: %v", i, err)
			continue
		}
		
		textContent.WriteString(text)
		textContent.WriteString("\n")
	}
	
	return &ExtractedContent{
		RawText:   strings.TrimSpace(textContent.String()),
		PageCount: numPages,
		Metadata: map[string]interface{}{
			"format":    "PDF",
			"pages":     numPages,
			"encrypted": isEncrypted,
		},
	}, nil
}

// extractFromDOCX extracts text from DOCX files
func (e *DocumentExtractor) extractFromDOCX(filePath string) (*ExtractedContent, error) {
	log.Printf("[EXTRACTOR] Extracting text from DOCX: %s", filePath)
	
	// Open DOCX file
	r, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DOCX file: %v", err)
	}
	defer r.Close()
	
	// Extract text content
	docx := r.Editable()
	text := docx.GetContent()
	
	// Clean up the text (remove excessive whitespace)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(text, "\n\n")
	text = strings.TrimSpace(text)
	
	return &ExtractedContent{
		RawText:   text,
		PageCount: 1, // DOCX doesn't have clear page boundaries in text extraction
		Metadata: map[string]interface{}{
			"format": "DOCX",
		},
	}, nil
}

// extractFromTXT extracts text from plain text files
func (e *DocumentExtractor) extractFromTXT(filePath string) (*ExtractedContent, error) {
	log.Printf("[EXTRACTOR] Reading text file: %s", filePath)
	
	// Open text file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open text file: %v", err)
	}
	defer file.Close()
	
	// Read all content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read text file: %v", err)
	}
	
	// Convert to string and clean up
	text := string(content)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.TrimSpace(text)
	
	return &ExtractedContent{
		RawText:   text,
		PageCount: 1,
		Metadata: map[string]interface{}{
			"format": "TXT",
		},
	}, nil
}

// ExtractClientData extracts structured data from raw text using patterns
func (e *DocumentExtractor) ExtractClientData(rawText string, patterns []ContentPattern) map[string]interface{} {
	log.Printf("[EXTRACTOR] Extracting structured data from %d characters of text", len(rawText))
	
	extractedData := make(map[string]interface{})
	
	for _, pattern := range patterns {
		matches := pattern.Pattern.FindStringSubmatch(rawText)
		if len(matches) > 1 {
			// Found a match
			value := matches[1] // First capture group
			if pattern.Processor != nil {
				value = pattern.Processor(value).(string)
			}
			extractedData[pattern.FieldName] = value
			log.Printf("[EXTRACTOR] Found %s: %s", pattern.FieldName, value)
		} else if pattern.Required {
			log.Printf("[EXTRACTOR] Warning: Required field %s not found", pattern.FieldName)
		}
	}
	
	return extractedData
}

// SanitizeText cleans and normalizes extracted text
func (e *DocumentExtractor) SanitizeText(text string) string {
	// Remove excessive whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	// Remove control characters
	text = regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`).ReplaceAllString(text, "")
	
	// Trim and normalize
	text = strings.TrimSpace(text)
	
	return text
}

// GetSupportedFormats returns the list of supported file formats
func (e *DocumentExtractor) GetSupportedFormats() []string {
	return e.SupportedFormats
}

// IsFormatSupported checks if a file format is supported
func (e *DocumentExtractor) IsFormatSupported(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, format := range e.SupportedFormats {
		if ext == format {
			return true
		}
	}
	return false
}