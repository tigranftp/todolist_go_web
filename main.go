package main

import (
	"fmt"
	"site/API"
)

func main() {
	siteAPI, err := API.NewSiteAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = siteAPI.Start()
	if err != nil {
		fmt.Println(err)
	}
}
