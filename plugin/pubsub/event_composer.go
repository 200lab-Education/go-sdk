package pb

import (
	"github.com/200Lab-Education/go-sdk/sdkcm"
)

type Opt func(*composer)

type composer struct {
	evt *Event
}

func EventComposer(title string, data interface{}, opts ...Opt) *composer {
	comp := new(composer)
	comp.evt = &Event{Title: title, Data: data}

	for _, o := range opts {
		o(comp)
	}

	return comp
}

func (c *composer) Event() *Event {
	return c.evt
}

func WithSenderIdAndObject(senderId, obj string) Opt {
	return func(c *composer) {
		if c.evt.Author == nil {
			c.evt.Author = &EntityDetail{Id: senderId, Object: obj}
			return
		}
		c.evt.Author.Id = senderId
		c.evt.Author.Object = obj
	}
}

func WithSenderNameAndImage(name string, image *sdkcm.Image) Opt {
	return func(c *composer) {
		if c.evt.Author == nil {
			c.evt.Author = &EntityDetail{Name: name, Image: image}
			return
		}
		c.evt.Author.Name = name
		c.evt.Author.Image = image
	}
}

func WithReceiverIdAndObject(receiverId, obj string) Opt {
	return func(c *composer) {
		if c.evt.Receiver == nil {
			c.evt.Receiver = &EntityDetail{Id: receiverId, Object: obj}
			return
		}
		c.evt.Receiver.Id = receiverId
		c.evt.Receiver.Object = obj
	}
}

func WithReceiverNameAndImage(name string, image *sdkcm.Image) Opt {
	return func(c *composer) {
		if c.evt.Receiver == nil {
			c.evt.Receiver = &EntityDetail{Name: name, Image: image}
			return
		}
		c.evt.Receiver.Name = name
		c.evt.Receiver.Image = image
	}
}
