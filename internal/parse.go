package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	Json     any             `json:"json,omitempty"`
}

type keyWithValues struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type UserResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func ParseJsonInput(file []byte) ([]UserRequest, error) {
	var newReq []UserRequest
	err := json.Unmarshal(file, &newReq)
	if err != nil {
		return newReq, fmt.Errorf("ParseJsonInput: unmarshal: %w", err)
	}

	return newReq, nil
}

func ConvertInputsToReqs(reqs []UserRequest) ([]Request, error) {
	if len(reqs) == 0 {
		return nil, errors.New("no requests provided")
	}
	newRequest := []Request{}
	for _, r := range reqs {
		conv, err := convertInputToReq(r)
		if err != nil {
			return nil, fmt.Errorf("ConvertInputsToReqs: %w", err)
		}
		newRequest = append(newRequest, *conv)
	}
	return newRequest, nil
}
func convertInputToReq(userReq UserRequest) (*Request, error) {
	jsonBody, err := json.Marshal(userReq.Body.Json)
	if err != nil {
		return nil, fmt.Errorf("convertInputToReq: error marshaling json field: %w", err)
	}
	return &Request{
		Name:    userReq.Name,
		Method:  userReq.HttpMethod,
		Url:     userReq.Url,
		Headers: convertMapStrings(userReq.Headers),
		Body: Body{
			Type:     userReq.Body.Type,
			FormData: convertMapStrings(userReq.Body.FormData),
			Raw:      userReq.Body.Raw,
			Binary:   userReq.Body.Binary,
			JsonData: jsonBody,
		},
		Params: Params{
			Values: convertMapStrings(userReq.Params),
		},
	}, nil
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

func ListRequests(rs []UserRequest) string {
	nameHeader := "Name"
	descHeader := "Description"
	methodHeader := "Method"
	urlHeader := "URL"
	nameMaxWidth := len(nameHeader)
	descMaxWidth := len(descHeader)
	methodMaxWidth := len(methodHeader)
	urlMaxWidth := len(urlHeader)
	rows := [][4]string{}
	rows = append(rows, [4]string{nameHeader, descHeader, methodHeader, urlHeader})
	for _, r := range rs {
		rows = append(rows, [4]string{r.Name, r.Description, r.HttpMethod, r.Url})
		if len(r.Name) > nameMaxWidth {
			nameMaxWidth = len(r.Name)
		}
		if len(r.Description) > descMaxWidth {
			descMaxWidth = len(r.Description)
		}
		if len(r.HttpMethod) > methodMaxWidth {
			methodMaxWidth = len(r.HttpMethod)
		}
		if len(r.Url) > urlMaxWidth {
			urlMaxWidth = len(r.Url)
		}
	}

	listStr := ""
	for _, row := range rows {
		rowStr := row[0] + strings.Repeat(" ", nameMaxWidth+2-len(row[0]))
		rowStr += row[1] + strings.Repeat(" ", descMaxWidth+2-len(row[1]))
		rowStr += row[2] + strings.Repeat(" ", methodMaxWidth+2-len(row[2]))
		rowStr += row[3] + strings.Repeat(" ", urlMaxWidth+2-len(row[3]))
		rowStr += "\n"
		listStr += rowStr
	}
	return listStr
}
