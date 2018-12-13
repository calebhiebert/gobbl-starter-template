package linkr

// Link represents an object that will be serialized to json, then to base64
type Link struct {
	// RedirectTo should be a URL that the user will be redirected to after clicking the link
	RedirectTo string `json:"r"`

	// PSID is the psid of the user this link was built for
	PSID string `json:"psid"`

	// Data contains any custom data that should be added to links
	Data CustomLinkData `json:"d"`
}

// CustomLinkData contains any data you want to associate with outgoing links
// TODO modify this to add custom data
type CustomLinkData struct {
	Metadata string `json:"met"`
}
