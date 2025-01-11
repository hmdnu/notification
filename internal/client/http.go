package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hmdnu/bot/internal/cookie"
	"github.com/hmdnu/bot/internal/utils"
)

type Headers struct {
	key   string
	value string
}

type Student struct {
	Nim      string
	Password string
}

type Options struct {
	headers        []Headers
	excludedCookie []string
}

var cookieJar *cookie.CookieJar

func NewHtppClient(jar *cookie.CookieJar) {
	cookieJar = jar
}

func login() error {
	data := Student{
		Nim:      utils.Env.Nim,
		Password: utils.Env.Password,
	}

	baseUrl := utils.Env.SiakadUrl
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

	res, err := fetch(http.MethodPost, baseUrl+"/login", payload, Options{headers: headers})
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	return nil
}

func serializeCookies(cookies []cookie.Cookie, excludedCookies []string) string {
	var serializedCookies []string

	for _, currentCookie := range cookies {
		if !includeString(excludedCookies, currentCookie.Key) {
			serializedCookies = append(serializedCookies, fmt.Sprintf("%s=%s", currentCookie.Key, currentCookie.Value))
		}
	}

	return strings.Join(serializedCookies, "; ")
}

func fetch(method string, path string, data io.Reader, options Options) (*http.Response, error) {
	var cookieString string

	rawCookies := cookieJar.Entries()
	cookieString = serializeCookies(rawCookies, options.excludedCookie)

	req, err := http.NewRequest(method, path, data)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Cookie", cookieString)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; GoClient/1.0)")

	if options.headers != nil {
		for _, header := range options.headers {
			req.Header.Set(header.key, header.value)
		}
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if cookies := res.Cookies(); cookies != nil {
		cookieJar.Store(cookies)
	}

	return res, nil
}

func CollectCookies() {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Collecting SIAKAD cookies", collectSiakadCookies},
		{"Attempting login", login},
		{"Collecting homepage cookies", collectHomepageCookies},
		{"Collecting SLC cookies", collectSlcCookies},
	}

	for _, step := range steps {
		fmt.Println(step.name)

		if err := step.fn(); err != nil {
			log.Fatalf("%s failed: %v", step.name, err)
		}

		fmt.Printf("%s success\n", step.name)
	}
}

func collectSiakadCookies() error {
	_, err := fetch(http.MethodHead, utils.Env.SiakadUrl, nil, Options{})
	return err
}

func collectHomepageCookies() error {
	_, err := fetch(http.MethodHead, utils.Env.SiakadUrl+"/beranda", nil, Options{})
	return err
}

func collectSlcCookies() error {
	res1, err1 := fetch(http.MethodHead, utils.Env.SlcUrl, nil, Options{})
	res2, err2 := fetch(http.MethodHead, utils.Env.SlcUrl, nil, Options{})

	if err1 != nil || err2 != nil {
		return fmt.Errorf("failed to fetch initial SLC cookies: %w", err1)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res1.Body)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res2.Body)

	return nil
}

// Check if slice of strings contains a value
func includeString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
