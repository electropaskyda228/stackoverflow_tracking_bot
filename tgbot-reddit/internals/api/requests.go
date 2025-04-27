package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	. "tgbot-reddit/internals/common"
)

type TrackInfo struct {
	UserId     uint
	QuestionId uint
}

var uri string = os.Getenv("LOCAL_HOST") + os.Getenv("STACKOVERFLOW_SERVER_PORT") + "/"

func CheckUserExisting(userId string) (bool, error) {
	client := http.Client{Timeout: 3 * time.Second}

	data := url.Values{}
	data.Add("user", userId)

	req, err := http.NewRequest(http.MethodGet, uri+"user/checking", nil)
	if err != nil {
		return false, err
	}
	req.URL.RawQuery = data.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, &BadRequest{}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body) == "true", nil
}

func MakeUser(userId string, name string) error {
	client := http.Client{Timeout: 3 * time.Second}

	data := url.Values{}
	data.Add("id", userId)
	data.Add("name", name)

	resp, err := client.PostForm(uri+"user/add", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return &FailedAddUser{}
	}
	return nil
}

func Listing(userId string) ([]string, error) {
	client := http.Client{Timeout: 3 * time.Second}

	data := url.Values{}
	data.Add("user", userId)

	req, err := http.NewRequest(http.MethodGet, uri+"question/list", nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = data.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, &BadRequest{}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []TrackInfo
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	answer := make([]string, 0)
	for _, track := range result {
		answer = append(answer, strconv.FormatUint(uint64(track.QuestionId), 10))
	}

	return answer, nil

}

func Tracking(userId string, questionId string) error {
	client := http.Client{Timeout: 3 * time.Second}

	data := url.Values{}
	data.Add("id", questionId)
	data.Add("user", userId)

	resp, err := client.PostForm(uri+"question/track", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return &FailedAddTracking{}
	}
	return nil
}

func Untracking(userId string, questionId string) error {
	client := http.Client{Timeout: 3 * time.Second}

	data := url.Values{}
	data.Add("ids", questionId)
	data.Add("user", userId)

	resp, err := client.PostForm(uri+"question/untrack", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return &FailedUntrack{}
	}
	return nil
}
