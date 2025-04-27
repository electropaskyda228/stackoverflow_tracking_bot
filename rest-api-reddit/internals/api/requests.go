package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	. "rest-api-reddit/internals/common"
	"time"
)

type ResponseQuestion struct {
	Items []Question
}

const uri = "https://api.stackexchange.com/2.3/questions/"

func CheckConstants() {
	if os.Getenv("STACKOVERFLOW_KEY") == "" {
		log.Println("Os's constants have not been found")
		os.Exit(1)
	}
}

func GetQuestion(ids string) (*Question, error) {
	client := http.Client{Timeout: 3 * time.Second}

	data := url.Values{}
	data.Add("site", "stackoverflow")
	data.Add("key", os.Getenv("STACKOVERFLOW_KEY"))

	req, err := http.NewRequest(http.MethodGet, uri+ids, nil)
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
	var question ResponseQuestion
	err = json.Unmarshal(body, &question)
	if err != nil {
		return nil, err
	}

	idsUint, err := StringToUint(ids)
	if err != nil {
		return nil, err
	}
	if len(question.Items) == 0 {
		return nil, &APIError{}
	}
	question.Items[0].ID = idsUint
	return &question.Items[0], nil

}
