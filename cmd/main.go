package main

import (
	"github.com/hmdnu/bot/internal/client"
	"github.com/hmdnu/bot/internal/collector"
	"github.com/hmdnu/bot/internal/cookie"
)

func main() {
	cookieJar := cookie.NewCookieJar()
	client.NewHtppClient(cookieJar)
	client.CollectCookies()
	collector.Collector()
}
