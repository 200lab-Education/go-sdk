package cloudinary

import "time"

type VideoResult struct {
	PublicID     string    `json:"public_id"`
	Version      int       `json:"version"`
	Signature    string    `json:"signature"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	Format       string    `json:"format"`
	ResourceType string    `json:"resource_type"`
	CreatedAt    time.Time `json:"created_at"`
	Tags         []string  `json:"tags"`
	Bytes        int       `json:"bytes"`
	Type         string    `json:"type"`
	Etag         string    `json:"etag"`
	URL          string    `json:"url"`
	SecureURL    string    `json:"secure_url"`
	Audio        Audio     `json:"audio"`
	Video        Video     `json:"video"`
	FrameRate    float64   `json:"frame_rate"`
	BitRate      int       `json:"bit_rate"`
	Duration     float64   `json:"duration"`
	Error        Error     `json:"error,omitempty"`
}
type Audio struct {
	Codec         string `json:"codec"`
	BitRate       string `json:"bit_rate"`
	Frequency     int    `json:"frequency"`
	Channels      int    `json:"channels"`
	ChannelLayout string `json:"channel_layout"`
}
type Video struct {
	PixFormat string `json:"pix_format"`
	Codec     string `json:"codec"`
	Level     int    `json:"level"`
	BitRate   string `json:"bit_rate"`
}

type Error struct {
	Message string `json:"message"`
}
