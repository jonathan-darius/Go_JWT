package model

import (
	"github.com/olahol/go-imageupload"
)

type Address struct {
	State    string `json:"state" bson:"state"`
	City     string `json:"city" bson:"city"`
	Postcode int    `json:"postcode" bson:"postcode"`
}

type User struct {
	Name    string  `json:"name" bson:"user_name"`
	Age     int     `json:"age" bson:"user_age"`
	Address Address `json:"address"`

	Path string
}

type User_IN struct {
	Name    string
	Age     int
	Address Address
	Path    string
}

type Image struct {
	ImageFile *imageupload.Image
}
