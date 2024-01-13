package export

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"ascii-art-web-export-file/Error"
)

func Export(w http.ResponseWriter, r *http.Request, docType, asciiResult string) {
	// Generate the appropriate file extension and content type based on the document type
	fileExt, contentType := getFileInfo(docType)

	// Generate a unique file name for the exported file
	fileName := generateFileName(fileExt)

	// Create the file in a temporary directory
	filePath := filepath.Join(os.TempDir(), fileName)
	file, err := os.Create(filePath)
	if err != nil {
		Error.RenderErrorPage(w,http.StatusInternalServerError , "Failed to create the exported file")
		return
	}
	defer file.Close()

	// Write the ASCII result to the file
	_, err = file.WriteString(asciiResult)
	if err != nil {
		Error.RenderErrorPage(w,http.StatusInternalServerError , "Failed to write to the exported file")
		return
	}

	// Set the response headers for the file download
	fileInfo, err := file.Stat()
	if err != nil {
		Error.RenderErrorPage(w,http.StatusInternalServerError , "Failed to get file information")
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	contentLength := strconv.FormatInt(fileInfo.Size(), 10)
	w.Header().Set("Content-Length", contentLength)

	// Serve the file for download
	http.ServeFile(w, r, filePath)

	// Remove the temporary file after download
	os.Remove(filePath)
}

func getFileInfo(docType string) (fileExt, contentType string) {
	switch docType {
	case "Plain Text":
		return ".txt", "text/plain"
	case "Rich Text Format":
		return ".rtf", "application/rtf"
	case "Markdown":
		return ".md", "text/markdown"
	case "Word":
		return ".doc", "application/msword"
	default:
		return ".txt", "text/plain"
	}
}

func generateFileName(fileExt string) string {
	return fmt.Sprintf("exported_result_%s%s", generateRandomString(8), fileExt)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomString)
}
