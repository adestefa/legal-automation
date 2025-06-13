package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ICloudDocument represents a document in iCloud
type ICloudDocument struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	Modified    time.Time `json:"modified"`
	Type        string    `json:"type"`
	IsDirectory bool      `json:"isDirectory"`
}

// ICloudService handles iCloud Drive operations
type ICloudService struct {
	testPath string
}

// NewICloudService creates a new iCloud service instance
func NewICloudService() *ICloudService {
	return &ICloudService{
		testPath: "/Users/corelogic/satori-dev/clients/proj-mallon/test_icloud",
	}
}

// GetRootFolders returns the top-level folders in iCloud Drive
func (s *ICloudService) GetRootFolders(username, appPassword string) ([]ICloudDocument, error) {
	// Always try test directory first for development
	var icloudPath string
	if _, err := os.Stat(s.testPath); err == nil {
		icloudPath = s.testPath
		fmt.Printf("Using test iCloud directory: %s\n", s.testPath)
	} else {
		// Get iCloud Drive path on macOS
		icloudPath = "/Users/" + getCurrentUser() + "/Library/Mobile Documents/com~apple~CloudDocs"
		
		// Check if iCloud Drive is available
		if _, err := os.Stat(icloudPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("iCloud Drive not available or not synced")
		}
	}
	
	// List directories in iCloud Drive root
	dirs, err := os.ReadDir(icloudPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read iCloud Drive: %v", err)
	}
	
	var folders []ICloudDocument
	for i, dir := range dirs {
		if dir.IsDir() && !strings.HasPrefix(dir.Name(), ".") {
			info, err := dir.Info()
			if err != nil {
				continue
			}
			
			folder := ICloudDocument{
				ID:          fmt.Sprintf("icloud_folder_%d", i),
				Name:        dir.Name(),
				Path:        "/" + dir.Name(),
				IsDirectory: true,
				Modified:    info.ModTime(),
				Size:        0, // Directories don't have size
			}
			folders = append(folders, folder)
		}
	}
	
	fmt.Printf("Found %d real iCloud folders\n", len(folders))
	return folders, nil
}

// GetSubfolders returns subfolders within a specific iCloud directory
func (s *ICloudService) GetSubfolders(username, appPassword, parentFolder string) ([]ICloudDocument, error) {
	// Always try test directory first for development
	var icloudPath string
	if _, err := os.Stat(s.testPath); err == nil {
		icloudPath = s.testPath
		fmt.Printf("Using test iCloud directory for subfolders: %s\n", s.testPath)
	} else {
		// Get iCloud Drive path on macOS
		icloudPath = "/Users/" + getCurrentUser() + "/Library/Mobile Documents/com~apple~CloudDocs"
		
		// Check if iCloud Drive is available
		if _, err := os.Stat(icloudPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("iCloud Drive not available or not synced")
		}
	}
	
	// Clean the parent folder path
	cleanParent := strings.TrimPrefix(parentFolder, "/")
	fullPath := filepath.Join(icloudPath, cleanParent)
	
	// Check if parent folder exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("folder does not exist: %s", parentFolder)
	}
	
	// List subdirectories
	dirs, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folder %s: %v", parentFolder, err)
	}
	
	var subfolders []ICloudDocument
	for i, dir := range dirs {
		if dir.IsDir() && !strings.HasPrefix(dir.Name(), ".") {
			info, err := dir.Info()
			if err != nil {
				continue
			}
			
			subfolder := ICloudDocument{
				ID:          fmt.Sprintf("icloud_subfolder_%d", i),
				Name:        dir.Name(),
				Path:        parentFolder + "/" + dir.Name(),
				IsDirectory: true,
				Modified:    info.ModTime(),
				Size:        0,
			}
			subfolders = append(subfolders, subfolder)
		}
	}
	
	fmt.Printf("Found %d subfolders in %s\n", len(subfolders), parentFolder)
	return subfolders, nil
}

// GetDocuments returns documents from a specific iCloud folder
func (s *ICloudService) GetDocuments(username, appPassword, folderPath string) ([]ICloudDocument, error) {
	// Always try test directory first for development
	var icloudPath string
	if _, err := os.Stat(s.testPath); err == nil {
		icloudPath = s.testPath
		fmt.Printf("[DEBUG] Using test iCloud directory for documents: %s\n", s.testPath)
	} else {
		// Get iCloud Drive path on macOS
		icloudPath = "/Users/" + getCurrentUser() + "/Library/Mobile Documents/com~apple~CloudDocs"
		
		// Check if iCloud Drive is available
		if _, err := os.Stat(icloudPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("iCloud Drive not available or not synced")
		}
		fmt.Printf("[DEBUG] Using real iCloud directory: %s\n", icloudPath)
	}
	
	// Clean the folder path
	cleanPath := strings.TrimPrefix(folderPath, "/")
	fullPath := filepath.Join(icloudPath, cleanPath)
	
	// Use root if path is empty
	if cleanPath == "" {
		fullPath = icloudPath
	}
	
	fmt.Printf("[DEBUG] Resolved full path: %s\n", fullPath)
	
	// Check if folder exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("[ERROR] Folder does not exist: %s (resolved to %s)\n", folderPath, fullPath)
		return nil, fmt.Errorf("folder does not exist: %s", folderPath)
	}
	
	// List all items in the folder
	items, err := os.ReadDir(fullPath)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read folder %s: %v\n", fullPath, err)
		return nil, fmt.Errorf("failed to read folder %s: %v", folderPath, err)
	}
	
	fmt.Printf("[DEBUG] Found %d items in directory %s\n", len(items), fullPath)
	
	var documents []ICloudDocument
	for i, item := range items {
		// Skip hidden files
		if strings.HasPrefix(item.Name(), ".") {
			fmt.Printf("[DEBUG] Skipping hidden file: %s\n", item.Name())
			continue
		}
		
		info, err := item.Info()
		if err != nil {
			fmt.Printf("[DEBUG] Error getting file info for %s: %v\n", item.Name(), err)
			continue
		}
		
		// Determine file type
		fileType := "unknown"
		if !item.IsDir() {
			ext := strings.ToLower(filepath.Ext(item.Name()))
			switch ext {
			case ".pdf":
				fileType = "pdf"
			case ".docx":
				fileType = "docx"
			case ".doc":
				fileType = "doc"
			case ".txt":
				fileType = "txt"
			case ".jpg", ".jpeg", ".png":
				fileType = "image"
			default:
				fileType = strings.TrimPrefix(ext, ".")
			}
		}
		
		doc := ICloudDocument{
			ID:          fmt.Sprintf("icloud_doc_%d", i),
			Name:        item.Name(),
			Path:        folderPath + "/" + item.Name(),
			IsDirectory: item.IsDir(),
			Modified:    info.ModTime(),
			Size:        info.Size(),
			Type:        fileType,
		}
		documents = append(documents, doc)
		fmt.Printf("[DEBUG] Added document: %s (%s, %d bytes)\n", doc.Name, doc.Type, doc.Size)
	}
	
	fmt.Printf("[SUCCESS] Found %d documents in iCloud folder %s\n", len(documents), folderPath)
	return documents, nil
}

// Helper function to get the current macOS username
func getCurrentUser() string {
	return "corelogic" // Simplified for development
}