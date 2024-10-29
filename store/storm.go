package store

import (
	"fmt"
	"time"
)

type URLMapping struct {
	ID          int       `storm:"id,increment"`
	UserId      string    `storm:"index"`
	ShortUrl    string    `storm:"index"`
	OriginalUrl string    `storm:"index"`
	CreatedAt   time.Time `storm:"index"`
}

func (s *StorageService) SaveToBolt(urlMapping URLMapping) error {
	err := s.storm.Save(&urlMapping)
	if err != nil {
		return fmt.Errorf("failed to save url mapping: %v", err)
	}
	return nil
}

func (s *StorageService) RetrieveFromBolt(shortUrl string) (URLMapping, error) {
	var urlMapping URLMapping
	err := s.storm.One("ShortUrl", shortUrl, &urlMapping)
	if err != nil {
		return URLMapping{}, fmt.Errorf("failed to retrieve original url: %v", err)
	}
	return urlMapping, nil
}
