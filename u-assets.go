package main

import (
	"fmt"
	"os"
)

// UStaticAssetURL will generate a URL to a given static asset
func UStaticAssetURL(name string) string {
	return fmt.Sprintf("%s/static/%s", os.Getenv("SELF_HOST"), name)
}
