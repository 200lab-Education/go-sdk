package imgprocessing

import (
	"flag"
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/200Lab-Education/go-sdk/sdkcm"
)

var (
	ErrImgProcessingHostMissing = sdkcm.CustomError("ErrImgProcessingHostMissing", "Img Processing service host is missing")
)

type imgProcessing struct {
	name   string
	prefix string
	logger logger.Logger
	cfg    imgProcessingConfig
}

type imgProcessingConfig struct {
	host string
}

func New(prefix ...string) *imgProcessing {
	pre := "img-processing"

	if len(prefix) > 0 {
		pre = prefix[0]
	}

	return &imgProcessing{
		name:   "img-processing",
		prefix: pre,
	}
}

func (imgproc *imgProcessing) Get() interface{} {
	return imgproc
}

func (imgproc *imgProcessing) Name() string {
	return imgproc.name
}

func (imgproc *imgProcessing) InitFlags() {
	flag.StringVar(&imgproc.cfg.host, fmt.Sprintf("%s-%s", imgproc.GetPrefix(), "host"), "", "img processing host")
}

func (imgproc *imgProcessing) GetPrefix() string {
	return imgproc.prefix
}

func (imgproc *imgProcessing) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (imgproc *imgProcessing) Configure() error {
	imgproc.logger = logger.GetCurrent().GetLogger(imgproc.Name())

	if err := imgproc.cfg.check(); err != nil {
		imgproc.logger.Errorln(err)
		return err
	}

	return nil
}

func (imgproc *imgProcessing) Run() error {
	return imgproc.Configure()
}

func (cfg *imgProcessingConfig) check() error {
	if len(cfg.host) < 0 {
		return ErrImgProcessingHostMissing
	}
	return nil
}
