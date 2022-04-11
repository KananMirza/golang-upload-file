package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func mainHandler(w http.ResponseWriter, _ *http.Request) {
	temp, err := template.ParseFiles("public/index.html")
	checkError(err)
	checkError(temp.Execute(w, nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	checkError(r.ParseMultipartForm(10 << 20))

	file, header, err := r.FormFile("file")
	checkError(err)

	defer func(file multipart.File) {
		err = file.Close()
		checkError(err)
	}(file)

	tempFile, err := ioutil.TempFile("public/", "upload-*"+filepath.Ext(header.Filename))
	checkError(err)

	defer func(tempFile *os.File) {
		err = tempFile.Close()
		checkError(err)
	}(tempFile)

	fileByte, err := ioutil.ReadAll(file)
	checkError(err)

	_, err = tempFile.Write(fileByte)
	checkError(err)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/upload", uploadHandler)
	fmt.Print("Server starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
