package urltool

import (
	"errors"
	"net/url"
	"path"
)

// GetBaeUrl 获取URL路径的最后一节
// singdile.space/mn -> mn
func GetBaseUrl(longurl string) (string, error) {
	urlmap, err := url.Parse(longurl)
	if err != nil {
		return "", err
	}
	if len(urlmap.Host) == 0 {
		return "", errors.New("need a valid url with host")
	}
	return path.Base(urlmap.Path), nil
}
