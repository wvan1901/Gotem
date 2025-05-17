package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type UserRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	HttpMethod  string          `json:"method"`
	Url         string          `json:"url"`
	Headers     []keyWithValues `json:"headers,omitempty"`
	Body        bodyInput       `json:"body"`
	Params      []keyWithValues `json:"parmeters,omitempty"`
}

type bodyInput struct {
	Type     string          `json:"type"`
	FormData []keyWithValues `json:"data,omitempty"`
	Raw      string          `json:"raw,omitempty"`
	Binary   []byte          `json:"binary,omitempty"` //This is base64 encoded string
}

type keyWithValues struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type UserResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func GetJsonFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("GetJsonFile: os: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			// TODO: Write to stderr
			fmt.Println(err)
		}
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("GetJsonFile: io: %w", err)
	}

	return fileBytes, nil
}

func ParseJsonInput(file []byte) ([]UserRequest, error) {
	var newReq []UserRequest
	err := json.Unmarshal(file, &newReq)
	if err != nil {
		return newReq, fmt.Errorf("ParseJsonInput: unmarshal: %w", err)
	}

	return newReq, nil
}

func ConvertInputsToReqs(reqs []UserRequest) []Request {
	if len(reqs) == 0 {
		return nil
	}
	newRequest := []Request{}
	for _, r := range reqs {
		newRequest = append(newRequest, convertInputToReq(r))
	}
	return newRequest
}
func convertInputToReq(userReq UserRequest) Request {
	return Request{
		Name:    userReq.Name,
		Method:  userReq.HttpMethod,
		Url:     userReq.Url,
		Headers: convertMapStrings(userReq.Headers),
		Body: Body{
			Type:     userReq.Body.Type,
			FormData: convertMapStrings(userReq.Body.FormData),
			Raw:      userReq.Body.Raw,
			Binary:   userReq.Body.Binary,
		},
		Params: Params{
			Values: convertMapStrings(userReq.Params),
		},
	}
}

func convertMapStrings(in []keyWithValues) map[string][]string {
	if len(in) == 0 {
		return nil
	}
	newMap := map[string][]string{}
	for _, keyWithVals := range in {
		newMap[keyWithVals.Key] = keyWithVals.Values
	}

	return newMap
}

func ConvertResponse(r Response) UserResponse {
	return UserResponse{
		StatusCode: r.StatusCode,
		Body:       string(r.Body),
	}
}

func (u *UserResponse) JsonString() (string, error) {
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		return "", fmt.Errorf("JsonString: %w", err)
	}
	return string(jsonBytes), nil
}
