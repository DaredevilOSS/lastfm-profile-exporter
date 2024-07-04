package internal

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const containerName = "debug"

var accountName = "lastfmprofileexports" // os.Getenv("ACCOUNT_NAME")
var accountKey = os.Getenv("ACCOUNT_KEY")

func uploadToBlobStorage(filename string) (*string, error) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	blobName := filepath.Base(filename)
	log.Printf("Upload file %s as blob %s\n", filename, blobName)
	azureUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)
	parsedUrl, err := url.Parse(azureUrl)
	if err != nil {
		return nil, err
	}

	containerURL := azblob.NewContainerURL(*parsedUrl, pipeline)
	blobURL := containerURL.NewBlockBlobURL(blobName)
	_, err = azblob.UploadFileToBlockBlob(context.Background(), file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   azblob.BlockBlobMaxUploadBlobBytes,
		Parallelism: 8,
	})
	if err != nil {
		return nil, err
	}
	blobUrlStr := blobURL.String()
	return &blobUrlStr, nil
}

func WriteToFile(username string, method string, contents []string) (*string, error) {
	uploadTimestampStr := time.Now().Format("1970_01_01T00_00_00")
	filename := fmt.Sprintf("/tmp/%s_%s_%s.ndjson", username, method, uploadTimestampStr)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(contents, "\n"))
	if err != nil {
		return nil, err
	}
	log.Printf("Wrote top tracks file of %d lines\n", len(contents))
	return &filename, nil
}

func upload(inputChannel chan string) {
	for filename := range inputChannel {
		_, err := uploadToBlobStorage(filename)
		if err != nil {
			log.Println(err)
		}
	}
}
