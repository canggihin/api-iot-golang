package service

import (
	"context"
	"errors"
	"mqtt-golang-rainfall-prediction/models"
	"mqtt-golang-rainfall-prediction/repository"
)

type Service interface {
	InsertData(ctx context.Context, data models.SensorData) error
}

type service struct {
	repositories repository.RepositoryInterface
}

func NewService(repositories repository.RepositoryInterface) *service {
	return &service{
		repositories: repositories,
	}
}

func (s *service) InsertData(ctx context.Context, data models.SensorData) error {
	err := s.repositories.InsertData(ctx, data)
	if err != nil {
		return errors.New("failed to insert data")
	}
	return nil
}
