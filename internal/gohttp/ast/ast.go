package ast

import (
	"net/http"
	"slices"
)

type Progarm struct {
	Requests []UserRequest
	// If we do global label here would be a good idea
	// GlobalLabels map[string]string
}

// Do we need labels to be global or local?
// For now we will have all labels be relative to the request, no global labels for now
type UserRequest struct {
	Name            string // Name of the request, requried
	Description     string // Request Description, not requried
	HttpRequestBody string // Is http body that follows http rfc
	HttpMethod      string
	Url             string
	ExtraLabels     map[string]string
}

func (u UserRequest) IsValid() (bool, map[string]string) {
	problems := map[string]string{}
	if u.Name == "" {
		problems["name"] = "is empty"
	}
	httpVerbs := []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}
	if !slices.Contains(httpVerbs, u.HttpMethod) {
		problems["http verb"] = "invalid"
	}
	if len(problems) == 0 {
		return true, nil
	}

	return false, problems
}
