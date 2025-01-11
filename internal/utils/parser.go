package utils

import "io"

func ParseToText(body io.ReadCloser) (string, error) {
	buffer, err := io.ReadAll(body)

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			return
		}
	}(body)

	return string(buffer), nil
}
