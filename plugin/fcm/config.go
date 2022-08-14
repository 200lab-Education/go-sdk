package fcm

import (
	"errors"
	"github.com/200Lab-Education/go-sdk/sdkcm"
)

var fcmErrors = map[string]string{
	"MissingRegistration": "Device is missing registration",
	"InvalidRegistration": "Device is invalid registration",
	"NotRegistered":       "Device is not registered",
	"InvalidPackageName":  "Make sure the message was addressed to a registration token whose package name matches the value passed in the request",
	"MismatchSenderId":    "A registration token is tied to a certain group of senders. When a client app registers for FCM, it must specify which senders are allowed to send messages",
	"InvalidParameters":   "Check that the provided parameters have the right name and type",

	"MessageTooBig": "Check that the total size of the payload data included in a message does not exceed FCM limits: 4096 bytes for most messages, or 2048 bytes in the case of " +
		"messages to topics. This includes both the keys and the values.",

	"InvalidDataKey": "Check that the payload data does not contain " +
		"a key (such as from, or gcm, or any value prefixed by google) that is used internally by FCM. Note that some " +
		"words (such as collapse_key) are also used by FCM but are allowed in the payload, in which case the payload " +
		"value will be overridden by the FCM value",

	"InvalidTtl":                "Check that the value used in time_to_live is an integer representing a duration in seconds between 0 and 2,419,200",
	"Unavailable":               "The server couldn't process the request in time",
	"InternalServerError":       "The server encountered an error while trying to process the request",
	"DeviceMessageRateExceeded": "The rate of messages to a particular device is too high. If an iOS app sends messages at a rate exceeding APNs limits, it may receive this error message",
	"TopicsMessageRateExceeded": "The rate of messages to subscribers to a particular topic is too high",
	"InvalidApnsCredential":     "A message targeted to an iOS device could not be sent because the required APNs authentication key was not uploaded or has expired",
}

var (
	ErrNotifyNotSuccess  = errors.New("can't send notification")
	ErrNotificationEmpty = errors.New("notification can't not be empty")

	ErrMissingRegistration       = createCustomError("MissingRegistration")
	ErrInvalidRegistration       = createCustomError("InvalidRegistration")
	ErrNotRegistered             = createCustomError("NotRegistered")
	ErrInvalidPackageName        = createCustomError("InvalidPackageName")
	ErrMismatchSenderId          = createCustomError("MismatchSenderId")
	ErrInvalidParameters         = createCustomError("InvalidParameters")
	ErrMessageTooBig             = createCustomError("MessageTooBig")
	ErrInvalidDataKey            = createCustomError("InvalidDataKey")
	ErrInvalidTtl                = createCustomError("InvalidTtl")
	ErrUnavailable               = createCustomError("Unavailable")
	ErrInternalServerError       = createCustomError("InternalServerError")
	ErrDeviceMessageRateExceeded = createCustomError("DeviceMessageRateExceeded")
	ErrTopicsMessageRateExceeded = createCustomError("TopicsMessageRateExceeded")
	ErrInvalidApnsCredential     = createCustomError("ErrInvalidApnsCredential")
)

func createCustomError(errKey string) sdkcm.ErrorWithKey {
	return sdkcm.CustomError(errKey, fcmErrors[errKey])
}
