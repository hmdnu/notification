package client

import (
	"net/http"

	"github.com/hmdnu/bot/config"
)

func FetchSubjectContent() (*http.Response, error) {
	headers := []Headers{
		{key: "Sec-Gpc", value: "1"},
		{key: "Sec-Fetch-Mode", value: "cors"},
		{key: "Sec-Fetch-Site", value: "same-origin"},
		{key: "Sec-Fetch-Dest", value: "empty"},
		{key: "Referer", value: config.Env.SlcUrl + "/spada"},
		{key: "Host", value: config.Env.SlcUrl},
		{key: "Pragma", value: "no-cache"},
	}

	res, err := fetch(http.MethodGet, config.Env.SlcUrl, nil, headers)

	if err != nil {
		return nil, err
	}

	return res, nil
}
