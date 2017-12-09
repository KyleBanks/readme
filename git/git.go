package git

import (
	"fmt"
	"net/url"
)

func RepositoryUrl(username, repository string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("https://github.com/%s/%s", username, repository))
}
