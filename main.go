package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/upload", handler)

	println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error while uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	var now = time.Now()
	var fileName = now.String() + header.Filename + ".pdf"

	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile("files", "template.docx")
	if err != nil {
		panic(err)
	}

	io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", "http://localhost:3000/forms/libreoffice/convert", buffer)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	resultFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer resultFile.Close()

	_, err = io.Copy(resultFile, response.Body)
	if err != nil {
		panic(err)
	}

	if err != nil {
		http.Error(w, "Error while sending the file to the API", http.StatusBadRequest)
		return
	}
	fileD, err := os.Open(fileName)
	streamPDFbytes, err := ioutil.ReadAll(fileD)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// stream straight to client(browser)
	w.Header().Add("Content-type", "application/pdf")
	w.Write(streamPDFbytes)
	os.Remove(fileName)
}
