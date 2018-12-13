/*
Package linkr contains utilities for encoding and decoding a Link object
in base64. These links allow embedding of arbitrary data in URLs coming
from the Messenger bot. This allows for more fine grained tracking.
The base URL is pulled from the SELF_HOST environment variable */
package linkr

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/matoous/go-nanoid"
)

// Encode will encode the provided link as base64.
// this base64 is meant to be appended to a linkr url
// these urls are used for tracking user bot interest
func Encode(link *Link) (string, error) {
	jsonBytes, err := json.Marshal(link)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// Decode will decode a base64 string into a Link
func Decode(b64 string) (Link, error) {
	var link Link

	jsonBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64))
	if err != nil {
		return link, err
	}

	err = json.Unmarshal(jsonBytes, &link)
	if err != nil {
		return link, err
	}

	return link, nil
}

// URL will encode a link object into a complete URL
func URL(link *Link) (string, error) {
	b64, err := Encode(link)
	if err != nil {
		return "", err
	}

	// Inject some random data into the url so that facebook can never cache it
	rand, err := gonanoid.Nanoid(4)
	if err != nil {
		return "", err
	}

	return os.Getenv("SELF_HOST") + "/linkr/" + rand + "/link?d=" + url.QueryEscape(b64), nil
}

// MustURL will encode a link object into a complete URL
// panics in case of error
func MustURL(link *Link) string {
	urlString, err := URL(link)
	if err != nil {
		panic(err)
	}

	return urlString
}
