package fcm

type Response struct {
	Success int
	Fail    int
	Results []Result
}

type Result struct {
	MessageID   *string `json:"message_id,omitempty"`
	Error       error   `json:"error,omitempty"`
	DeviceToken string  `json:"device_token"`
}
