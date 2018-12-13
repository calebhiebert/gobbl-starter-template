package main

import (
	"github.com/calebhiebert/gobbl/messenger"
)

// USetMessengerProfile will set the facebook messenger profile variable
func USetMessengerProfile(messengerAPI *fb.MessengerAPI) (interface{}, error) {
	res, err := messengerAPI.MessengerProfile(&fb.MessengerProfile{
		GetStarted: fb.GetStarted{
			Payload: "GET_STARTED",
		},
	})
	return res, err
}
