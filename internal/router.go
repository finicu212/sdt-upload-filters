package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 64MB RAM, rest temp files
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// files uploaded to the form
	fileHeaders := r.MultipartForm.File["file"]
	for _, fileHeader := range fileHeaders {
		lf, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed opening file header: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		err = os.MkdirAll("./uploads/user", os.ModePerm)
		rf, err := os.Create(fmt.Sprintf("uploads/%s/%s", "user", fileHeader.Filename))
		defer rf.Close()

		_, err = io.Copy(rf, lf)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed writing to file: %s", err.Error()), http.StatusInternalServerError)
			return
		}

	}
}
