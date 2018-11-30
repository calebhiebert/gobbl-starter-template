package main

import (
	"github.com/calebhiebert/gobbl"
	"github.com/calebhiebert/gobbl-localization"
	"github.com/calebhiebert/gobbl/messenger"
)

// HGetStarted is for when the user first speaks to the bot
func HGetStarted(c *gbl.Context) {
	l := glocalize.GetCurrentLocalization(c)
	r := fb.CreateResponse(c)

	// Greeting text
	r.Text(l.T("welcome"))
}

// HDefaultFallback is for when the bot doesn't know what is going on
func HDefaultFallback(c *gbl.Context) {
	l := glocalize.GetCurrentLocalization(c)
	r := fb.CreateResponse(c)

	r.Text(l.T("fallback"))
}
