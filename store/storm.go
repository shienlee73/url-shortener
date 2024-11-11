package store

import (
	"fmt"
)

func (s *StorageService) CreateURLMapping(urlMapping URLMapping) error {
	err := s.storm.Save(&urlMapping)
	if err != nil {
		return fmt.Errorf("failed to save url mapping: %v", err)
	}
	return nil
}

func (s *StorageService) RetrieveURLMapping(shortUrl string) (URLMapping, error) {
	var urlMapping URLMapping
	err := s.storm.One("ShortUrl", shortUrl, &urlMapping)
	if err != nil {
		return URLMapping{}, fmt.Errorf("failed to retrieve original url: %v", err)
	}
	return urlMapping, nil
}

func (s *StorageService) CreateUser(user User) error {
	err := s.storm.Save(&user)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}
	return nil
}

func (s *StorageService) RetrieveUser(username string) (User, error) {
	var user User
	err := s.storm.One("Username", username, &user)
	if err != nil {
		return User{}, fmt.Errorf("failed to retrieve user: %v", err)
	}
	return user, nil
}

func (s *StorageService) CreateSession(session Session) error {
	err := s.storm.Save(&session)
	if err != nil {
		return fmt.Errorf("failed to save session: %v", err)
	}
	return nil
}

func (s *StorageService) RetrieveSession(sessionID string) (Session, error) {
	var session Session
	err := s.storm.One("ID", sessionID, &session)
	if err != nil {
		return Session{}, fmt.Errorf("failed to retrieve session: %v", err)
	}
	return session, nil
}

func (s *StorageService) DeleteSession(sessionID string) error {
	err := s.storm.DeleteStruct(&Session{ID: sessionID})
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}
	return nil
}
