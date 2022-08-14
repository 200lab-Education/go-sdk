package fcm

import gofcm "github.com/NaySoftware/go-fcm"

type Notification struct {
	Payload          map[string]string `json:"-"`
	CollapseKey      *string           `json:"-"`
	TimeToLive       int               `json:"-"`
	DelayWhileIdle   bool              `json:"-"`
	Title            string            `json:"title,omitempty"`
	Body             string            `json:"body,omitempty"`
	Icon             string            `json:"icon,omitempty"`
	Sound            string            `json:"sound,omitempty"`
	Badge            string            `json:"badge,omitempty"`
	Tag              string            `json:"tag,omitempty"`
	Color            string            `json:"color,omitempty"`
	ClickAction      string            `json:"click_action,omitempty"`
	BodyLocKey       string            `json:"body_loc_key,omitempty"`
	BodyLocArgs      string            `json:"body_loc_args,omitempty"`
	TitleLocKey      string            `json:"title_loc_key,omitempty"`
	TitleLocArgs     string            `json:"title_loc_args,omitempty"`
	AndroidChannelID string            `json:"android_channel_id,omitempty"`
}

func NewNotification(title string, opts ...Opt) *Notification {
	n := &Notification{Title: title}
	for _, o := range opts {
		o(n)
	}
	return n
}

func WithBody(body string) Opt {
	return func(n *Notification) {
		n.Body = body
	}
}

func WithIcon(icon string) Opt {
	return func(n *Notification) {
		n.Icon = icon
	}
}

func WithColor(color string) Opt {
	return func(n *Notification) {
		n.Color = color
	}
}

func WithSound(sound string) Opt {
	return func(n *Notification) {
		n.Sound = sound
	}
}

func WithBadge(badge string) Opt {
	return func(n *Notification) {
		n.Badge = badge
	}
}

func WithAction(action string) Opt {
	return func(n *Notification) {
		n.ClickAction = action
	}
}

func WithTag(tag string) Opt {
	return func(n *Notification) {
		n.Tag = tag
	}
}

func WithPayload(payload map[string]string) Opt {
	return func(n *Notification) {
		n.Payload = payload
	}
}

func WithCollapseKey(collapseKey string) Opt {
	return func(n *Notification) {
		if len(collapseKey) < 1 {
			return
		}
		n.CollapseKey = &collapseKey
	}
}

func (s *Notification) toNotificationPayload() *gofcm.NotificationPayload {
	return &gofcm.NotificationPayload{
		Title:            s.Title,
		Body:             s.Body,
		TitleLocKey:      s.TitleLocKey,
		TitleLocArgs:     s.TitleLocArgs,
		Tag:              s.Tag,
		Sound:            s.Sound,
		Icon:             s.Icon,
		Color:            s.Color,
		ClickAction:      s.ClickAction,
		BodyLocKey:       s.BodyLocKey,
		BodyLocArgs:      s.BodyLocArgs,
		Badge:            s.Badge,
		AndroidChannelID: s.AndroidChannelID,
	}
}

type Opt func(*Notification)
