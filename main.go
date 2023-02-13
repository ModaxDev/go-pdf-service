package main

import (
	"bytes"
	"github.com/joho/godotenv"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	start := time.Now()

	response, err := postFile(os.Getenv("URL_LIBREOFFICE"), "template.docx")
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	resultFile, err := os.Create("result.pdf")
	if err != nil {
		panic(err)
	}
	defer resultFile.Close()

	_, err = io.Copy(resultFile, response.Body)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	println("Temps d'exécution :", elapsed.Seconds(), "secondes")
}

func postFile(url string, filePath string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", filePath)
	if err != nil {
		return nil, err
	}

	io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	return client.Do(request)
}
