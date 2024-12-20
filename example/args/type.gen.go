// Code generated by 'xmodel gen' command.
// DO NOT EDIT!

package args

import (
	"github.com/swxctx/xmodel/mongo"
)

// User user info
type User struct {
	Id        int64  `key:"pri" json:"id"`
	Name      string `key:"uni" json:"name"`
	Age       int32  `json:"age"`
	UpdatedAt int64  `json:"updated_at"`
	CreatedAt int64  `json:"created_at"`
	DeletedTs int64  `json:"deleted_ts"`
}

// Meta comment...
type Meta struct {
	Id        mongo.ObjectId `json:"_id" bson:"_id" key:"pri"`
	Hobby     []string       `json:"hobby" bson:"hobby"`
	Tags      []string       `json:"tags" bson:"tags"`
	UpdatedAt int64          `json:"updated_at" bson:"updated_at"`
	CreatedAt int64          `json:"created_at" bson:"created_at"`
	DeletedTs int64          `json:"deleted_ts" bson:"deleted_ts"`
}