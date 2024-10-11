package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type FileMetadata struct {
	FileName string `json:"fileName"`
	FileID   string `json:"fileId"`
	FileSize int64  `json:"fileSize"`
}

var (
	jsonMutex sync.Mutex
)

func updateJSON(buffFileMetadata FileMetadata) error {
	jsonMutex.Lock()
	defer jsonMutex.Unlock()

	serverJSONFilePath := "./data/data.json"
	file, err := os.ReadFile(serverJSONFilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading JSON file: %w", err)
	}

	var jsonFileBuffer []FileMetadata
	if len(file) > 0 {
		if err := json.Unmarshal(file, &jsonFileBuffer); err != nil {
			return fmt.Errorf("error unmarshaling JSON: %w", err)
		}
	}

	jsonFileBuffer = append(jsonFileBuffer, buffFileMetadata)

	fileJSONBuff, err := json.MarshalIndent(jsonFileBuffer, "", "   ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(serverJSONFilePath, fileJSONBuff, 0644); err != nil {
		return fmt.Errorf("error writing JSON file: %w", err)
	}

	return nil
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Error processing file", http.StatusBadRequest)
		return
	}

	fileMap := r.MultipartForm.File["file"]
	for _, fileHeader := range fileMap {
		serverFile, err := os.Create("./uploads/" + fileHeader.Filename)
		if err != nil {
			log.Printf("Error creating file: %v", err)
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer serverFile.Close()

		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error opening uploaded file: %v", err)
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		_, err = io.Copy(serverFile, file)
		if err != nil {
			log.Printf("Error copying file: %v", err)
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		err = updateJSON(FileMetadata{
			FileName: fileHeader.Filename,
			FileID:   "0", // Consider generating a unique ID
			FileSize: fileHeader.Size,
		})
		if err != nil {
			log.Printf("Error updating JSON: %v", err)
			http.Error(w, "Error updating file metadata", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully")
}
