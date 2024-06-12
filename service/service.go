package service

import (
	"context"
	"errors"
	"mqtt-golang-rainfall-prediction/models"
	"mqtt-golang-rainfall-prediction/repository"
	"sort"
	"time"
)

type Service interface {
	InsertData(ctx context.Context, data models.SensorData) error
	GetData(ctx context.Context) ([]models.SensorData, error)
	GetDataByDay(ctx context.Context) ([]models.SensorDataByDay, error)
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

func (s *service) GetData(ctx context.Context) ([]models.SensorData, error) {
	data, err := s.repositories.GetData(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if data == nil {
		return nil, errors.New("data not found")
	}

	// Sort data by FormattedTime, parsing the time to ensure correct ordering
	sort.Slice(data, func(i, j int) bool {
		// Parse the times
		t1, err1 := time.Parse("02-01-2006 15:04:05", data[i].FormattedTime)
		t2, err2 := time.Parse("02-01-2006 15:04:05", data[j].FormattedTime)
		if err1 != nil || err2 != nil {
			return false // Handle parse error, maybe log or handle differently
		}
		return t1.After(t2) // Use After to sort descending
	})

	return data, nil
}

func (s *service) GetDataByDay(ctx context.Context) ([]models.SensorDataByDay, error) {
	data, err := s.repositories.GetDataPerDay(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if data == nil {
		return nil, errors.New("data not found")
	}

	// Sort data by FormattedTime, parsing the time to ensure correct ordering
	sort.Slice(data, func(i, j int) bool {
		// Parse the times
		t1, err1 := time.Parse("02-01-2006", data[i].FormattedTime)
		t2, err2 := time.Parse("02-01-2006", data[j].FormattedTime)
		if err1 != nil || err2 != nil {
			return false // Handle parse error, maybe log or handle differently
		}
		return t1.After(t2) // Use After to sort descending
	})

	return data, nil
}
