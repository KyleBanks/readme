package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpResolver struct {
	baseUrl *url.URL
}

func NewGitHubHttpResolver(username, repository string) (*HttpResolver, error) {
	url, err := url.Parse(fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/", username, repository))
	if err != nil {
		return nil, err
	}

	return &HttpResolver{
		baseUrl: url,
	}, nil
}

// Resolve downloads a file relative to the GitHub repository. If the filename provided is an absolute URL,
// as opposed to relative, it will still be downloaded.
func (h *HttpResolver) Resolve(filename string) (string, error) {
	// Images can be `/example.png` but its not actually a root URL its relative to user/repo/master/
	if strings.HasPrefix(filename, "/") {
		filename = "." + filename
	}

	fileUrl, err := url.Parse(filename)
	if err != nil {
		return "", err
	}

	url := h.baseUrl.ResolveReference(fileUrl)
	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New(string(contents))
	}

	return string(contents), nil
}
