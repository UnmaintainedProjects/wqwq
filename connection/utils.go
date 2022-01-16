package connection

import (
	"crypto/rand"
	"encoding/json"
	"strconv"
)

func isRequest(line []byte) (Request, bool) {
	request := Request{}
	err := json.Unmarshal(line, &request)
	if err != nil {
		return request, false
	}
	if request.Method == "" {
		return request, false
	}
	return request, true
}

func isResponse(line []byte) (Response, bool) {
	response := Response{}
	err := json.Unmarshal(line, &response)
	if err != nil {
		return response, false
	}
	if response.Result == nil {
		return response, false
	}
	return response, true
}

func getRandomId() (string, error) {
	s := ""
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return s, err
	}
	for _, i := range b {
		s += strconv.Itoa(int(i))
	}
	return s, nil
}
