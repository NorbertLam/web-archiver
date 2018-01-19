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

// A Page represents a webpage holding its url and its html layout.
type Page struct {
	Url  string `json:"url,omitempty"`
	Html string `json:"html,omitempty"`
}

// A Key represents private information to be accessed.
type Key struct {
	AccountName   string
	AccountKey    string
	Url           string
	ContainerName string
}

// Env holds a Key accessed from environment variables.
type Env struct {
	eKey Key
}

var (
	blobCli storage.BlobStorageClient
)

// Create page is a POST request which passes an Env and takes in data.
// Data is read and saved in a Page struct.
// BlockBlob operations are conducted, creating an empty BlockBlob.
// Page data is uploaded to BlockBlob using PutBlock.
// PutBlockList then commits the new BlockBlob to the database.
func (env *Env) createPage(w http.ResponseWriter, req *http.Request) {
	var page Page
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&page); err != nil {
		log.Fatal(err)
	}

	pageTitle := page.Url + ".html"
	data := page.Html
	length := len(data)

	byteArray := make([]byte, length)
	for i := 0; i < length; i++ {
		byteArray[i] = byte(data[i])
	}

	container := blobCli.GetContainerReference(env.eKey.ContainerName)
	blobRef := container.GetBlobReference(pageTitle)
	blobRef.CreateBlockBlob(nil)

	blockID := base64.StdEncoding.EncodeToString([]byte(pageTitle))
	if err := blobRef.PutBlock(blockID, []byte(byteArray), nil); err != nil {
		fmt.Printf("put block failed: %v", err)
	}

	list, err := blobRef.GetBlockList(storage.BlockListTypeUncommitted, nil)
	if err != nil {
		fmt.Printf("get block list failed: %v", err)
	}

	uncommittedBlocksList := make([]storage.Block, len(list.UncommittedBlocks))
	for i := range list.UncommittedBlocks {
		uncommittedBlocksList[i].ID = list.UncommittedBlocks[i].Name
		uncommittedBlocksList[i].Status = storage.BlockStatusUncommitted
	}

	if err = blobRef.PutBlockList(uncommittedBlocksList, nil); err != nil {
		fmt.Printf("put block list failed: %v", err)
	}

}

// Sets up env variables and opens connection to Azure blob storage client.
// Initializes router with a handler for creating and uploading a html file to the Azure database.
func main() {
	var key Key
	if err := envconfig.Process("myserver", &key); err != nil {
		fmt.Println(err)
	}
	env := &Env{eKey: key}

	client, err := storage.NewBasicClient(env.eKey.AccountName, env.eKey.AccountKey)
	if err != nil {
		fmt.Println(err)
	}
	blobCli = client.GetBlobService()

	router := mux.NewRouter()
	router.HandleFunc("/page/", env.createPage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
