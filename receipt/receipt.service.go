package receipt

import (
	"encoding/json"
	"fmt"
	"github.com/goWebServices/one-o-one/cors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const receiptPath = "receipts"

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		receiptList, err := GetReceipts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(receiptList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		r.ParseMultipartForm(5 << 20) // 5mb goes in memory. If exceeds, data goes on on disk
		file, handler, err := r.FormFile("receipt")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer file.Close()
		f, err := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

//SetupRoutes
func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.Middleware(receiptHandler))
}


