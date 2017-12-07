package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"text/template"
)

const (
	githubHttpTemplate = "https://raw.githubusercontent.com/{{.Username}}/{{.Repository}}/master/{{.Filename}}"
)

func NewGitHubHttpResolver() *HttpResolver {
	return &HttpResolver{
		template: template.Must(template.New("github").Parse(githubHttpTemplate)),
	}
}

type HttpResolver struct {
	template *template.Template
}

func (h *HttpResolver) Resolve(username, repository, filename string) (string, error) {
	var url bytes.Buffer
	err := h.template.Execute(&url, map[string]string{
		"Username":   username,
		"Repository": repository,
		"Filename":   filename,
	})
	if err != nil {
		return "", err
	}

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
