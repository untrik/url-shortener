package savehelper

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/untrik/url-shortener/internal/lib/random"
)

const maxRetries = 5

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func SaveWithRetry(urlSaver URLSaver, reqURL string, aliasLength int) (int64, string, error) {
	for i := 0; i <= maxRetries; i++ {
		alias := random.NewRandomString(aliasLength)
		id, err := urlSaver.SaveURL(reqURL, alias)
		if err != nil {
			return id, alias, nil
		}
		if aliasUniqueErr, ok := err.(*pq.Error); ok && aliasUniqueErr.Code == "23505" {
			continue
		}
		return 0, "", err
	}
	return 0, "", fmt.Errorf("failed to generate unique alias after %d attempts", maxRetries)
}
