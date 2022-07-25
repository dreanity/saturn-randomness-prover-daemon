package drand

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type Round struct {
	Round             uint64 `json:"round"`
	Randomness        string `json:"randomness"`
	Signature         string `json:"signature"`
	PreviousSignature string `json:"previous_signature"`
}

func GetRound(urls []string, roundNumber uint64) (*Round, error) {
	urlsLen := len(urls)
	urlIndex := rand.Uint32() % uint32(urlsLen)

	round := Round{}
	requestEndpoint := fmt.Sprint(urls[urlIndex], "/", "/public/", roundNumber)
	resp, err := http.Get(requestEndpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &round); err != nil {
		return nil, err
	}

	return &round, nil
}
