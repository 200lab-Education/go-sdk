package main

import (
	"context"
	// "context"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/fcm"
	"log"
)

var (
	deviceToken = "dn2k1g_Cutc:APA91bG88Pfb2Y7bhmifBIhRTm1p0uVcHBtDwa3LRTyLUq2vdguNC2ORk5BYvog1zsIHxZ5WDSBx3po7wU6eUXGfncYkeKoXx6sqld3JX_YbLAXEc1tL4P55a0MNi4GCJgnMXCLZclPv"
)

func main() {
	service := goservice.New(
		goservice.WithName("demo"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(fcm.New("fcm")),
	)
	_ = service.Init()

	client := service.MustGet("fcm").(fcm.FirebaseCloudMessaging)
	client.ShowPrintResult(true)

	collapseKey := "messages is my name"

	notification := fcm.NewNotification("title",
		fcm.WithIcon("ic_notification"),
		fcm.WithColor("#18d821"),
		fcm.WithSound("default"),
		fcm.WithCollapseKey(collapseKey),
		fcm.WithTag(collapseKey),
	)

	sendToDeviceList(client, notification)
	// sendToDevice(client, notification)
	// sendToTopic(client, notification)
}

func sendToDevice(client fcm.FirebaseCloudMessaging, notification *fcm.Notification) {
	res, err := client.SendToDevice(context.Background(), deviceToken, notification)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
}

func sendToDeviceList(client fcm.FirebaseCloudMessaging, notification *fcm.Notification) {
	res, err := client.SendToDevices(context.Background(), []string{"xxx", deviceToken, "123"}, notification)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
}

func sendToTopic(client fcm.FirebaseCloudMessaging, notification *fcm.Notification) {
	res, err := client.SendToTopic(context.Background(), "/topics/chat", notification)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
}
