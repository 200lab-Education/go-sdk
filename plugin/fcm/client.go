package fcm

import (
	"context"
	gofcm "github.com/NaySoftware/go-fcm"
	"strconv"
)

func (s *fcmClient) ShowPrintResult(show bool) {
	s.showPrintResult = show
}

// Send notification to a topic
func (s *fcmClient) SendToTopic(ctx context.Context, topic string, notification *Notification) (*Response, error) {
	s.client.NewFcmTopicMsg(topic, notification.Payload)

	s.prepare(notification)

	status, err := s.send()
	if err != nil {
		return nil, err
	}
	response := s.parseResponse(status)

	return response, nil
}

// Send notification to single device
func (s *fcmClient) SendToDevice(ctx context.Context, deviceId string, notification *Notification) (*Response, error) {
	if notification == nil {
		return nil, ErrNotificationEmpty
	}

	s.client.NewFcmMsgTo(deviceId, notification.Payload)
	if ck := notification.CollapseKey; ck != nil {
		s.client.SetCollapseKey(*ck)
	}

	s.client.SetNotificationPayload(notification.toNotificationPayload())

	s.prepare(notification)

	status, err := s.send()
	if err != nil {
		return nil, err
	}
	response := s.parseResponse(status)

	return response, nil
}

// Send notification to multi-devices
func (s *fcmClient) SendToDevices(ctx context.Context, deviceIds []string, notification *Notification) (*Response, error) {
	if notification == nil {
		return nil, ErrNotificationEmpty
	}

	s.client.NewFcmRegIdsMsg(deviceIds, notification.Payload)
	s.client.SetNotificationPayload(notification.toNotificationPayload())
	if ck := notification.CollapseKey; ck != nil {
		s.client.SetCollapseKey(*ck)
	}

	s.prepare(notification)

	status, err := s.send()
	if err != nil {
		return nil, err
	}
	response := s.parseResponse(status)

	return response, nil
}

// Parse from NaySoftware response status -> fcm response
func (s *fcmClient) parseResponse(status *gofcm.FcmResponseStatus) *Response {
	response := &Response{
		Success: status.Success,
		Fail:    status.Fail,
		Results: make([]Result, len(status.Results)),
	}

	// case send TOPIC successfully
	if status.MsgId > 0 {
		response.Results = make([]Result, 1)
		msgId := strconv.FormatInt(status.MsgId, 10)
		response.Results[0] = Result{
			MessageID:   &msgId,
			DeviceToken: s.client.Message.To,
		}
		return response
	}

	for index, item := range status.Results {
		messageId := item["message_id"]
		err := item["error"]

		deviceToken := s.client.Message.To
		if len(deviceToken) < 1 && len(s.client.Message.RegistrationIds) > 0 {
			deviceToken = s.client.Message.RegistrationIds[index]
		}

		if len(messageId) > 0 {
			response.Results[index] = Result{
				MessageID:   &messageId,
				DeviceToken: deviceToken,
			}
		}

		if len(err) > 0 {
			response.Results[index] = Result{
				DeviceToken: deviceToken,
				Error:       createCustomError(err),
			}
		}
	}
	return response
}

// Prepare before sending message
func (s *fcmClient) prepare(notification *Notification) {
	if notification.CollapseKey != nil {
		s.client.SetCollapseKey(*notification.CollapseKey)
	}

	s.client.SetTimeToLive(notification.TimeToLive)
	s.client.SetDelayWhileIdle(notification.DelayWhileIdle)
}

// Send notification
func (s *fcmClient) send() (*gofcm.FcmResponseStatus, error) {
	status, err := s.client.Send()
	if err != nil {
		return nil, err
	}

	if s.showPrintResult {
		status.PrintResults()
	}

	return status, nil
}
