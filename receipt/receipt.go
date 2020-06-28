package receipt

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ReceiptDirectory string = filepath.Join("uploads")

type Receipt struct {
	ReceiptName string `json:"name"`
	UploadDate time.Time `json:"uploadDate"`
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptPath))
	if len(urlPathSegment[1:]) > 1 {
		log.Println("Invalid URL while trying to download")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := urlPathSegment[1:][0]
	file, err := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if err != nil {
		log.Printf("Requested File Not found! [%s]", fileName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	fContentType := http.DetectContentType(fHeader)
	stat, err := file.Stat()
	if err != nil {
		log.Printf("Error occurred while checking the stats for file: [%s]", fileName)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fSize := 


}

func GetReceipts() ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	files, err := ioutil.ReadDir(ReceiptDirectory)
	if err != nil {
		log.Println("Error Occurred while reading 'uploads' directory")
		return nil, err
	}
	for _, f := range files {
		receipts = append(receipts, Receipt{ReceiptName: f.Name(), UploadDate: f.ModTime()})
	}

	return receipts, nil
}