package sdkcm

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// For reading
type SQLModel struct {
	// Real id in db, we would't show it
	ID uint32 `json:"-" gorm:"id,PRIMARY_KEY"`
	// Fake id, we will public it
	FakeID    UID       `json:"id" gorm:"-"`
	Status    *int      `json:"status,omitempty" gorm:"column:status;default:1;"`
	CreatedAt *JSONTime `json:"created_at,omitempty;" gorm:"column:created_at;"`
	UpdatedAt *JSONTime `json:"updated_at,omitempty;" gorm:"column:updated_at;"`
}

func (sm *SQLModel) GenUID(objType int, shardID uint32) *SQLModel {
	sm.FakeID = NewUID(sm.ID, objType, shardID)
	return sm
}

func NewSQLModelWithStatus(status int) *SQLModel {
	t := JSONTime(time.Now().UTC())
	return &SQLModel{
		Status:    &status,
		CreatedAt: &t,
		UpdatedAt: &t,
	}
}

func (sm *SQLModel) ToID() *SQLModel {
	sm.ID = sm.FakeID.localID
	return sm
}

// For creating
type SQLModelCreate struct {
	// Real id in db, we would't show it
	ID uint32 `json:"-" gorm:"id,PRIMARY_KEY"`
	// Fake id, we will public it
	FakeID    UID       `json:"id" gorm:"-"`
	Status    int       `json:"status,omitempty" gorm:"column:status;default:1;"`
	CreatedAt *JSONTime `json:"created_at,omitempty;" gorm:"column:created_at;"`
	UpdatedAt *JSONTime `json:"updated_at,omitempty;" gorm:"column:updated_at;"`
}

func (sm *SQLModelCreate) GenUID(objType int, shardID uint32) {
	sm.FakeID = NewUID(sm.ID, objType, shardID)
}

func NewSQLModelCreateWithStatus(status int) SQLModelCreate {
	t := JSONTime(time.Now().UTC())
	return SQLModelCreate{
		Status:    status,
		CreatedAt: &t,
		UpdatedAt: &t,
	}
}

// Set time format layout. Default: 2006-01-02
func SetDateFormat(layout string) {
	dateFmt = layout
}

type JSON []byte

// This method for mapping JSON to json data type in sql
func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

// This method for scanning JSON from json data type in sql
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	if s, ok := value.([]byte); ok {
		*j = append((*j)[0:0], s...)
		return nil
	}

	return errors.New("invalid Scan Source")
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return nil, errors.New("object json is nil")
	}

	return *j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("object json is nil")
	}

	*j = JSON(data)
	return nil
}

type Image struct {
	ID            uint32 `json:"img_id,omitempty" bson:"img_id,omitempty"`
	FakeID        *UID   `json:"id,omitempty" bson:"-"`
	Url           string `json:"url" bson:"url"`
	FileName      string `json:"file_name,omitempty" bson:"file_name,omitempty"`
	OriginWidth   int    `json:"org_width" bson:"org_width"`
	OriginHeight  int    `json:"org_height" bson:"org_height"`
	OriginUrl     string `json:"org_url" bson:"org_url"`
	CloudName     string `json:"cloud_name,omitempty" bson:"cloud_name"`
	CloudId       string `json:"cloud_id,omitempty" bson:"cloud_id"`
	DominantColor string `json:"dominant_color" bson:"dominant_color"`
	RequestId     string `json:"request_id,omitempty" bson:"-"`
	FileSize      uint32 `json:"file_size,omitempty" bson:"-"`
}

func (i *Image) HideSomeInfo() *Image {
	if i != nil {
		//i.CloudID = ""
		i.CloudId = ""
	}

	return i
}

func (i *Image) Fulfill(domain string) {
	i.Url = fmt.Sprintf("%s%s", domain, i.CloudId)
}

// This method for mapping Image to json data type in sql
func (i *Image) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	b, err := json.Marshal(i)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Image from date data type in sql
func (i *Image) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	if err := json.Unmarshal(v, i); err != nil {
		return err
	}
	return nil
}

type Images []Image

// This method for mapping Images to json array data type in sql
func (is *Images) Value() (driver.Value, error) {
	if is == nil {
		return nil, nil
	}

	b, err := json.Marshal(is)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Images from json array type in sql
func (is *Images) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	var imgs []Image

	if err := json.Unmarshal(v, &imgs); err != nil {
		return err
	}

	*is = Images(imgs)
	return nil
}

type Video struct {
	ID            uint32    `json:"img_id,omitempty" bson:"img_id,omitempty"`
	FakeID        *UID      `json:"id,omitempty" bson:"-"`
	Url           string    `json:"url" bson:"url"`
	OriginWidth   int       `json:"org_width" bson:"org_width"`
	OriginHeight  int       `json:"org_height" bson:"org_height"`
	OriginUrl     string    `json:"org_url" bson:"org_url"`
	CloudName     string    `json:"cloud_name,omitempty" bson:"cloud_name"`
	CloudId       string    `json:"cloud_id,omitempty" bson:"cloud_id"`
	DominantColor string    `json:"dominant_color" bson:"dominant_color"`
	RequestId     string    `json:"request_id,omitempty" bson:"-"`
	FileSize      uint32    `json:"file_size,omitempty" bson:"-"`
	Format        string    `json:"format,omitempty"`
	Audio         AudioInfo `json:"audio"`
	Video         VideoInfo `json:"video"`
	FrameRate     float64   `json:"frame_rate"`
	BitRate       int       `json:"bit_rate"`
	Duration      float64   `json:"duration"`
}

type AudioInfo struct {
	Codec         string `json:"codec" bson:"codec,omitempty"`
	BitRate       string `json:"bit_rate" bson:"bit_rate,omitempty"`
	Frequency     int    `json:"frequency" bson:"frequency,omitempty"`
	Channels      int    `json:"channels" bson:"channels,omitempty"`
	ChannelLayout string `json:"channel_layout" bson:"channel_layout,omitempty"`
}
type VideoInfo struct {
	PixFormat string `json:"pix_format" bson:"pix_format,omitempty"`
	Codec     string `json:"codec" bson:"codec,omitempty"`
	Level     int    `json:"level" bson:"level,omitempty"`
	BitRate   string `json:"bit_rate" bson:"bit_rate,omitempty"`
}

func (i *Video) HideSomeInfo() *Video {
	if i != nil {
		//i.CloudID = ""
		i.CloudId = ""
	}
	return i
}

func (i *Video) Fulfill(domain string) {
	i.Url = fmt.Sprintf("%s%s", domain, i.CloudId)
}

// This method for mapping Video to json data type in sql
func (i *Video) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	b, err := json.Marshal(i)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Video from date data type in sql
func (i *Video) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	if err := json.Unmarshal(v, i); err != nil {
		return err
	}
	return nil
}

type Videos []Video

// This method for mapping Videos to json array data type in sql
func (is *Videos) Value() (driver.Value, error) {
	if is == nil {
		return nil, nil
	}

	b, err := json.Marshal(is)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Videos from json array type in sql
func (is *Videos) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	var videos []Video

	if err := json.Unmarshal(v, &videos); err != nil {
		return err
	}

	*is = Videos(videos)
	return nil
}

type Media struct {
	ID            uint32     `json:"img_id,omitempty" bson:"img_id,omitempty"`
	FakeID        *UID       `json:"id,omitempty" bson:"-"`
	Type          string     `json:"type" bson:"type,omitempty"`
	Url           string     `json:"url" bson:"url,omitempty"`
	OriginWidth   int        `json:"org_width" bson:"org_width,omitempty"`
	OriginHeight  int        `json:"org_height" bson:"org_height,omitempty"`
	OriginUrl     string     `json:"org_url,omitempty" bson:"org_url,omitempty"`
	FileName      string     `json:"file_name,omitempty" bson:"file_name,omitempty"`
	CloudName     string     `json:"cloud_name,omitempty" bson:"cloud_name,omitempty"`
	CloudId       string     `json:"cloud_id,omitempty" bson:"cloud_id,omitempty"`
	DominantColor string     `json:"dominant_color,omitempty" bson:"dominant_color,omitempty"`
	RequestId     string     `json:"request_id,omitempty" bson:"-"`
	FileSize      uint32     `json:"file_size,omitempty" bson:"-"`
	Format        string     `json:"format,omitempty" bson:"format,omitempty"`
	Thumbnail     *Image     `json:"thumbnail,omitempty" bson:"thumbnail,omitempty"`
	Audio         *AudioInfo `json:"audio,omitempty" bson:"audio,omitempty"`
	Video         *VideoInfo `json:"video,omitempty" bson:"video,omitempty"`
	FrameRate     float64    `json:"frame_rate,omitempty" bson:"frame_rate,omitempty"`
	BitRate       int        `json:"bit_rate,omitempty" bson:"bit_rate,omitempty"`
	Duration      float64    `json:"duration,omitempty" bson:"duration,omitempty"`
}

func (i *Media) Fulfill(domain string) {
	i.Url = fmt.Sprintf("%s%s", domain, i.CloudId)
}

// This method for mapping Media to json data type in sql
func (i *Media) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	b, err := json.Marshal(i)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Media from date data type in sql
func (i *Media) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	if err := json.Unmarshal(v, i); err != nil {
		return err
	}
	return nil
}

type Medias []Media

func (i *Medias) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	b, err := json.Marshal(i)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (i *Medias) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	if err := json.Unmarshal(v, i); err != nil {
		return err
	}
	return nil
}
