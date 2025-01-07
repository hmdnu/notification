package main

import (
	"fmt"
	"io"
	"log"

	"github.com/hmdnu/bot/client"
)

func main() {
	client.CollectCookies()

	res, err := client.FetchSubjectContent()

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	fmt.Println(string(body))

}
