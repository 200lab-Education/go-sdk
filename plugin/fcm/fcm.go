/*----------------------------------------------------------------*\
 * @author          Ly Nam <lyquocnam@live.com>
 * @copyright       2019 Viet Tran <viettranx@gmail.com>
 * @license         Apache-2.0
 * @description		Plugin to work with Firebase Cloud Messaging
 *----------------------------------------------------------------*/
package fcm

import (
	"context"
	"flag"
	"fmt"
	"github.com/200Lab-Education/go-sdk/logger"
	gofcm "github.com/NaySoftware/go-fcm"
)

type FirebaseCloudMessaging interface {
	// Send notification to a topic
	// @topic: /topics/chat
	SendToTopic(ctx context.Context, topic string, notification *Notification) (*Response, error)

	// Send notification to a device
	// @deviceId: fCZ4_yRHP5U:APA91bHJTY...
	SendToDevice(ctx context.Context, deviceId string, notification *Notification) (*Response, error)

	// Send notification to many devices
	// @deviceIds: []string{ "fCZ4_yRHP5U:APA91bHJTY..." }
	SendToDevices(ctx context.Context, deviceIds []string, notification *Notification) (*Response, error)

	// Show the response result of notification, use for debugging
	ShowPrintResult(show bool)

	// Get API Key
	APIKey() string
}

type fcmClient struct {
	name            string
	apiKey          string
	showPrintResult bool
	client          *gofcm.FcmClient
	logger          logger.Logger
}

func (s *fcmClient) GetPrefix() string {
	return s.Name()
}

func (s *fcmClient) Get() interface{} {
	return s
}

func New(name string) *fcmClient {
	return &fcmClient{
		name:            name,
		showPrintResult: false,
	}
}

func (s *fcmClient) APIKey() string {
	return s.apiKey
}

func (s *fcmClient) Name() string {
	return s.name
}

func (s *fcmClient) InitFlags() {
	flag.StringVar(&s.apiKey, fmt.Sprintf("%s-api-key", s.Name()), "", "firebase cloud messaging api key")
}

func (s *fcmClient) Configure() error {
	s.logger = logger.GetCurrent().GetLogger(s.Name())
	s.client = gofcm.NewFcmClient(s.apiKey)
	return nil
}

func (s *fcmClient) Run() error {
	return s.Configure()
}

func (s *fcmClient) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()
	return c
}
