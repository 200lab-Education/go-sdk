package sdkcm

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type MgoModel struct {
	PK        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Status    int           `json:"status" bson:"status,omitempty"`
	CreatedAt *JSONTime     `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt *JSONTime     `json:"updated_at" bson:"updated_at,omitempty"`
	DeletedAt *JSONTime     `json:"deleted_at" bson:"deleted_at,omitempty"`
}

func (md *MgoModel) PrepareForInsert(status int) {
	jsTime := JSONTime(time.Now().UTC())
	md.PK = bson.NewObjectId()
	md.CreatedAt = &jsTime
	md.UpdatedAt = &jsTime
	md.Status = status
}
