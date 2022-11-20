package internal

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"sdt-upload-filters/pkg/connection"
	"sdt-upload-filters/pkg/connection/pool"
	"sdt-upload-filters/pkg/handler"
)

// Use templates to feed html file
func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

var (
	conn connection.IConnection
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var err error
	conn, err = pool.Instance().GetConnection(context.TODO())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html")
	vars := map[string]interface{}{"UUID": conn.GetUUID()}
	//http.ServeFile(w, r, "submit.html")
	outputHTML(w, "submit.html", vars)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
		remoteFilePath := fmt.Sprintf("uploads/%s/%s", conn.GetUUID(), fileHeader.Filename)
		err = os.MkdirAll(path.Dir(remoteFilePath), os.ModePerm)
		rf, err := os.Create(remoteFilePath)
		defer rf.Close()

		chain := handler.StringSliceToChain(r.MultipartForm.Value["handlers"])
		chain.Handle(rf, lf)
		//_, err = io.Copy(rf, lf)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed writing to file: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		fmt.Printf("handled a file\n")
	}
	fmt.Printf("done with conn\n")
	err := pool.Instance().ReleaseConnection(conn)
	if err != nil {
		return
	}
	fmt.Printf("Released connection... going to serve login page\n")
	http.ServeFile(w, r, "index.html")
	//outputHTML(w, "index.html", vars)
}
