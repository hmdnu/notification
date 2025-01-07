package cookie

import (
	"net/http"
)

type Cookie struct {
	Key      string
	Value    string
	Path     string
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
}

type CookieJar struct {
	storage map[string]Cookie
}

func NewCookieJar() *CookieJar {
	return &CookieJar{
		storage: make(map[string]Cookie),
	}
}

func (jar *CookieJar) Store(rawCookies []*http.Cookie) {
	cookies := Parse(rawCookies)

	for _, cookie := range cookies {

		if cookie.Key == "siakad_sess" {
			jar.storage[cookie.Key] = cookie
		}
	}
}

func (jar *CookieJar) Get(name string) Cookie {
	return jar.storage[name]
}

func (jar *CookieJar) Entries() []Cookie {
	cookies := []Cookie{}

	for _, cookie := range jar.storage {
		cookies = append(cookies, cookie)
	}

	return cookies
}

func Parse(rawCookies []*http.Cookie) []Cookie {

	cookies := []Cookie{}

	for _, cookie := range rawCookies {
		cookies = append(cookies, Cookie{
			Key:      cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			HttpOnly: cookie.HttpOnly,
			Secure:   cookie.Secure,
			SameSite: cookie.SameSite,
		})
	}

	return cookies
}
