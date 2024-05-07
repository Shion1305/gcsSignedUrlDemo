package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func generateSignedURL(bucket, object string, credsFilePath string) string {
	client, err := storage.NewClient(
		context.Background(),
		option.WithCredentialsFile(credsFilePath))
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// set 15 minutes expiration
	expires := time.Now().Add(15 * time.Minute)

	url, err := client.Bucket(bucket).SignedURL(object,
		&storage.SignedURLOptions{
			Method:  "GET",
			Expires: expires,
			// add custom header
			Headers: []string{"test: aaa"},
			Scheme:  storage.SigningSchemeV4,
		})
	if err != nil {
		log.Fatalf("Unable to generate signed URL: %v", err)
	}
	return url
}

func DownloadFromUrl(url string) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// add custom header specified in generateSignedURL
	req.Header.Add("test", "aaa")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status code:", resp.StatusCode)
	// save to file
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("output/out.zip", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	bucket := "ynufes-mypage-staging-bucket"
	object := "test"
	credsFilePath := "./ynufes-mypage-staging-dev.json"

	signedURL := generateSignedURL(bucket, object, credsFilePath)
	fmt.Println("Signed URL:", signedURL)
	DownloadFromUrl(signedURL)
}
