package scrap

import (
	"errors"
	"io"
	"net/http"
	"unicode"
)

func CreateHttpRes(url string) (*http.Response, error) {
	res, err := http.Get(url)
  	
	if err != nil {
		return nil, err
  	}
  	defer res.Body.Close()

  	if res.StatusCode != 200 {
		return nil, err
	}

	return res, nil
}

func CreateHttpReq(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil) 
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, HandleHttpStatusErr(res)
	}
	return res.Body, nil

	// data, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	return "", err
	// }

	// return string(data), nil
}

func CreateHttpReqWithReferer(url, refererLink string) (string, error) {
	req, err := http.NewRequest("GET", url, nil) 
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req.Header.Add("Referer", refererLink)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", HandleHttpStatusErr(res)
	}
	// return res.Body, nil

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func HandleHttpStatusErr(res *http.Response) (error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return  err
	}

	return errors.New("status code error! "+ string(data) )
}


func Preprocessing(input string) string {
	var result string

	isSpace := false
	for _, char := range input {
		if char == ' '{
			if isSpace {
				continue
			} else {
				isSpace = true
			}
		} else {
			isSpace = false
		}
		
		if unicode.IsGraphic(char) {
			result += string(char)
		}
	}

	return result
} 
