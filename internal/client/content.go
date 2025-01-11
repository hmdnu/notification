package client

import (
	"net/http"

	"github.com/hmdnu/bot/internal/utils"
)

var hasMoodleSession = false

func FetchSubjectContent() (string, error) {

	headers := []Headers{
		{key: "Sec-Gpc", value: "1"},
		{key: "Sec-Fetch-Mode", value: "cors"},
		{key: "Sec-Fetch-Site", value: "same-origin"},
		{key: "Sec-Fetch-Dest", value: "empty"},
		{key: "Referer", value: utils.Env.SlcUrl + "/spada"},
		{key: "Host", value: utils.Env.SlcUrl},
		{key: "Pragma", value: "no-cache"},
	}

	res, err := fetch(http.MethodGet, utils.Env.SlcUrl, nil, Options{headers: headers})

	if err != nil {
		return "", err
	}

	str, err := utils.ParseToText(res.Body)

	if err != nil {
		return "", err
	}

	return str, nil
}

func FetchLmsContent(courseUrl string) (string, error) {
	_, err := fetch(http.MethodGet, utils.Env.SlcUrl+"/spada?gotourl="+courseUrl, nil, Options{})

	if err != nil {
		return "", err
	}

	if !hasMoodleSession {
		headers := []Headers{
			{key: "Sec-Fetch-Site", value: "same-site"},
			{key: "Referer", value: utils.Env.SlcUrl},
		}

		_, err := fetch(http.MethodGet, courseUrl, nil, Options{headers: headers})

		if err != nil {
			return "", err
		}

		if cookieJar.Has("MoodleSession") {
			hasMoodleSession = true
		}

	}

	res, err := fetch(http.MethodGet, courseUrl, nil, Options{})

	if err != nil {
		return "", err
	}

	str, err := utils.ParseToText(res.Body)

	if err != nil {
		return "", err
	}

	return str, nil
}
