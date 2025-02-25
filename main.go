package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Bucket string
	region   string
)

func main() {
	s3Bucket = os.Getenv("S3_BUCKET")
	region = os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("S3_ENDPOINT")

	if s3Bucket == "" || region == "" || accessKey == "" || secretKey == "" || endpoint == "" {
		log.Fatal("Missing required environment variables: S3_BUCKET, AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, S3_ENDPOINT")
	}

	http.HandleFunc("/download", downloadHandler)
	port := ":8080"
	fmt.Printf("Server started at %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileKey := r.URL.Query().Get("file")
	if fileKey == "" {
		http.Error(w, "file parameter is required", http.StatusBadRequest)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Credentials:      credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true), // Отключение SSL (если нужно)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create AWS session: %v", err), http.StatusInternalServerError)
		return
	}

	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileKey),
	}
	result, err := svc.GetObject(input)
	if err != nil {
		if reqErr, ok := err.(awserr.RequestFailure); ok && reqErr.StatusCode() == 404 {
			http.Error(w, fmt.Sprintf("File not found: %s", fileKey), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("failed to get object from S3: %v", err), http.StatusInternalServerError)
		return
	}
	defer result.Body.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileKey))
	w.Header().Set("Content-Type", *result.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))

	_, err = io.Copy(w, result.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to stream file: %v", err), http.StatusInternalServerError)
		return
	}
}
