package pb

import (
	"context"
	"fmt"
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Channel string

type Provider interface {
	Publish(ctx context.Context, channel Channel, data *Event) error
	Subscribe(ctx context.Context, channel Channel) (c <-chan *Event, close func())
}

type Entity interface {
	EntityID() string
	EntityName() string
	EntityObject() string
	EntityImage() *sdkcm.Image
}

type EntityDetail struct {
	Id     string       `json:"id"`
	Name   string       `json:"name"`
	Object string       `json:"object"`
	Image  *sdkcm.Image `json:"image"`
}

func (e *EntityDetail) GetId() string          { return e.Id }
func (e *EntityDetail) GetName() string        { return e.Name }
func (e *EntityDetail) GetObject() string      { return e.Object }
func (e *EntityDetail) GetImage() *sdkcm.Image { return e.Image }

type Event struct {
	Id         string        `json:"id"`
	Title      string        `json:"title"`
	Author     *EntityDetail `json:"author"`
	Receiver   *EntityDetail `json:"receiver"`
	Channel    Channel       `json:"channel"`
	Data       interface{}   `json:"data"`
	Ack        func()
	CreatedAt  time.Time `json:"created_at"`
	RemoteData []byte    `json:"remote_data"`
}

func (e Event) String() string {
	if e.Author != nil && e.Receiver != nil {
		return fmt.Sprintf("Event '%s' | Author: %s - ID: %s (%s) | Receiver: %s - ID: %s (%s)",
			e.Title,
			e.Author.Name,
			e.Author.Id,
			e.Author.Object,
			e.Receiver.Name,
			e.Receiver.Id,
			e.Receiver.Object,
		)
	}

	if e.Author != nil {
		return fmt.Sprintf("Event '%s' | Author: System | Receiver: %s - ID: %s (%s)",
			e.Title,
			e.Receiver.Name,
			e.Receiver.Id,
			e.Receiver.Object,
		)
	}

	return fmt.Sprintf("Event '%s' | Data: %v", e.Title, e.Data)
}

func (e *Event) GetID() string              { return e.Id }
func (e *Event) GetTitle() string           { return e.Title }
func (e *Event) GetAuthor() *EntityDetail   { return e.Author }
func (e *Event) GetReceiver() *EntityDetail { return e.Receiver }
func (e *Event) GetChannel() Channel        { return e.Channel }
func (e *Event) GetData() interface{}       { return e.Data }
func (e *Event) GetWhen() interface{}       { return e.CreatedAt }
func (e *Event) DoAck()                     { e.Ack() }
func (e *Event) SetChannel(c Channel)       { e.Channel = c }
func (e *Event) SetAck(f func())            { e.Ack = f }

func NewEvent(title string, author, receiver Entity, data interface{}) *Event {
	return &Event{
		Id:        bson.NewObjectId().Hex(),
		Title:     title,
		Author:    toEntity(author),
		Receiver:  toEntity(receiver),
		Data:      data,
		CreatedAt: time.Now().UTC(),
	}
}

func toEntity(e Entity) *EntityDetail {
	if e == nil {
		return nil
	}
	return &EntityDetail{e.EntityID(), e.EntityName(), e.EntityObject(), e.EntityImage()}
}
