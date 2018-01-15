package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Page struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type Key struct {
	AccountName   string
	AccountKey    string
	Url           string
	ContainerName string
}

var (
	blobCli storage.BlobStorageClient
)

func createPage(w http.ResponseWriter, req *http.Request) {
	var key Key
	envconfig.Process("myserver", &key)

	var p Page
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}

	pageTitle := p.Title + ".html"
	data := p.Body
	lent := len(data)

	s := make([]byte, lent)
	for i := 0; i < lent; i++ {
		s[i] = byte(data[i])
	}

	cnt := blobCli.GetContainerReference(key.ContainerName)
	b := cnt.GetBlobReference(pageTitle)
	b.CreateBlockBlob(nil)

	blockID := base64.StdEncoding.EncodeToString([]byte(pageTitle))
	err = b.PutBlock(blockID, []byte(s), nil)
	if err != nil {
		fmt.Println("put block failed: %v", err)
	}

	list, err := b.GetBlockList(storage.BlockListTypeUncommitted, nil)
	if err != nil {
		fmt.Println("get block list failed: %v", err)
	}

	uncommittedBlocksList := make([]storage.Block, len(list.UncommittedBlocks))
	for i := range list.UncommittedBlocks {
		uncommittedBlocksList[i].ID = list.UncommittedBlocks[i].Name
		uncommittedBlocksList[i].Status = storage.BlockStatusUncommitted
	}

	err = b.PutBlockList(uncommittedBlocksList, nil)
	if err != nil {
		fmt.Println("put block list failed: %v", err)
	}

}

func main() {
	var key Key
	err := envconfig.Process("myserver", &key)
	if err != nil {
		fmt.Println(err)
	}

	client, err := storage.NewBasicClient(key.AccountName, key.AccountKey)
	if err != nil {
		fmt.Println(err)
	}
	blobCli = client.GetBlobService()

	router := mux.NewRouter()
	router.HandleFunc("/page/", createPage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
