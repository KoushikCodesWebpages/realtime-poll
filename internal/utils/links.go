package utils

import (
	"realtime-poll/internal/models"

	"net/url"
	"time"
	"fmt"
) 

func AttachShareLink(p *models.Poll, baseURL string) {
	if p.Distribution.ShareID != "" {
		link := baseURL + "/poll/share/" + p.Distribution.ShareID
		p.Distribution.ShareID = link
	}
}

func BuildPageLink(basePath string, q url.Values, t time.Time, dir string) string {
	q.Set("cursor", t.Format(time.RFC3339))
	q.Set("dir", dir)

	return fmt.Sprintf("%s?%s", basePath, q.Encode())
}