package main

import (
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/imgprocessing"
	"github.com/200Lab-Education/go-sdk/sdkcm"
)

func main() {
	service := goservice.New(
		goservice.WithName("demo"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(imgprocessing.New("imgproc")),
	)
	_ = service.Init()

	imgproc := service.MustGet("imgproc").(imgprocessing.ImgProcessing)
	img, err := imgproc.ResizeFile("test.png", "scale", 0, 90)

	if err != nil {
		fmt.Printf("err: %+v", err.(sdkcm.AppError).Log)
	}
	fmt.Printf("%+v", img)
}
