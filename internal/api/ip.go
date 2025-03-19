package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const endpoint = "http://ip-api.com/json/"

var (
	ErrRequest       = errors.New("failed API request")
	ErrPrivateRange  = errors.New("private range")
	ErrReservedRange = errors.New("reserved range")
	ErrInvalidQuery  = errors.New("invalid query")
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Country string `json:"country"`
	Region  string `json:"regionName"`
	City    string `json:"city"`
}

func FindLocation(IP string) (r Response, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s?fields=status,message,country,regionName,city", endpoint, IP),
		nil,
	)
	if err != nil {
		return Response{}, fmt.Errorf("create HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return Response{}, fmt.Errorf("decode JSON: %w", err)
	}

	if r.Status == "fail" {
		switch r.Message {
		case "private range":
			return Response{}, ErrPrivateRange
		case "reserved range":
			return Response{}, ErrReservedRange
		case "invalid query":
			return Response{}, ErrInvalidQuery
		default:
			return Response{}, fmt.Errorf("%w: %s", ErrRequest, r.Message)
		}
	}

	return
}
