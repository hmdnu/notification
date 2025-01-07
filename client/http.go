package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	config "github.com/hmdnu/bot/config"
	"github.com/hmdnu/bot/cookie"
)

type Headers struct {
	key   string
	value string
}

type Student struct {
	Nim      string
	Password string
}

var cookieJar = cookie.NewCookieJar()

func login() error {
	data := Student{
		Nim:      config.Env.Nim,
		Password: config.Env.Password,
	}

	baseUrl := config.Env.SiakadUrl

	formData := url.Values{}

	formData.Set("username", data.Nim)
	formData.Set("password", data.Password)

	payload := bytes.NewBufferString(formData.Encode())

	headers := []Headers{
		{key: "Content-Type", value: "application/x-www-form-urlencoded"},
		{key: "Referer", value: baseUrl + "/login/index/err/6"},
		{key: "Accept", value: "application/json"},
		{key: "X-Requested-With", value: "XMLHttpRequest"},
	}

	res, err := fetch("POST", baseUrl, payload, headers)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	return nil
}

func fetch(method string, path string, data io.Reader, headers []Headers) (*http.Response, error) {
	req, err := http.NewRequest(method, path, data)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	)

	// iteration headers
	for _, header := range headers {
		if headers != nil {
			req.Header.Set(header.key, header.value)
		}
	}

	cookies := res.Cookies()

	if cookies != nil {
		cookieJar.Store(cookies)
	} else {
		cookieJar.Store(nil)
	}

	return res, nil
}

func collectSiakadCookies() error {
	_, err := fetch(http.MethodHead, config.Env.SiakadUrl, nil, nil)

	if err != nil {
		return err
	}

	return nil
}

func collectHomepageCookies() error {
	_, err := fetch(http.MethodHead, config.Env.SiakadUrl+"/beranda", nil, nil)

	if err != nil {
		return err
	}

	return nil
}

func collectSlcCookies() error {
	_, err := fetch(http.MethodHead, config.Env.SlcUrl, nil, nil)
	_, err2 := fetch(http.MethodHead, config.Env.SlcUrl, nil, nil)

	if err != nil && err2 != nil {
		return err
	}

	return nil
}

func CollectCookies() {
	// =====================================
	fmt.Println("collecting siakad cookies")
	err := collectSiakadCookies()

	if err != nil {
		log.Fatal("failed collecting cookies ", err.Error())
	}
	fmt.Println("collecting siakad cookies success")

	// =====================================
	fmt.Println("attempting login")
	err = login()

	if err != nil {
		log.Fatal("failed to login ", err.Error())
	}
	fmt.Println("login success")

	// =====================================
	fmt.Println("collecting siakad homepage cookies")
	err = collectHomepageCookies()

	if err != nil {
		log.Fatal("failed collecting siakad homepage cookies ", err.Error())
	}
	fmt.Println("success collecting siakad homepage cookies")

	// =====================================
	fmt.Println("collecting slc cookies")
	err = collectSlcCookies()

	if err != nil {
		log.Fatal("failed collecting slc cookies ", err.Error())
	}
	fmt.Println("success collecting slc cookies")
	// =====================================
}
