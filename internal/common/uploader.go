package common

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const containerName = "debug"

var accountName = "lastfmprofileexports" // os.Getenv("ACCOUNT_NAME")
var accountKey = os.Getenv("ACCOUNT_KEY")

func uploadToBlobStorage(filename string) error {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return fmt.Errorf("failed to create credentials: %w", err)
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	blobName := filepath.Base(filename)
	log.Printf("Upload file %s as blob %s\n", filename, blobName)
	azureUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)
	parsedUrl, err := url.Parse(azureUrl)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	containerURL := azblob.NewContainerURL(*parsedUrl, pipeline)
	blobURL := containerURL.NewBlockBlobURL(blobName)
	_, err = azblob.UploadFileToBlockBlob(context.Background(), file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   azblob.BlockBlobMaxUploadBlobBytes,
		Parallelism: 16,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func WriteToFile(username string, method string, contents []string) string {
	filename := fmt.Sprintf("/tmp/%s_%s_%s.ndjson", username, method, randString(16))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(contents, "\n"))
	if err != nil {
		log.Fatalf("failed to write to file: %s", err)
	}
	log.Printf("Wrote top tracks file of %d lines\n", len(contents))
	return filename
}

func Uploader(inputChannel FileChannel) {
	for filename := range inputChannel {
		err := uploadToBlobStorage(filename)
		if err != nil {
			log.Println(err)
		}
	}
}
