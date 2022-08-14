package main

import (
	"context"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/aws"
	"log"
)

func main() {
	service := goservice.New(
		goservice.WithName("demo"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(aws.New("aws")),
	)
	_ = service.Init()
	logoFile := "logo.png" // put this file on project root to test

	s3 := service.MustGet("aws").(aws.S3)
	url, err := s3.Upload(context.Background(), logoFile, "media")
	if err != nil {
		log.Fatalln(err)
	}

	urls := []string{"media/1572142633410254000.png", "media/1572143325272547000.png"} // put fileKey need remove into array

	if err := s3.DeleteImages(context.Background(), urls); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(url)
}
