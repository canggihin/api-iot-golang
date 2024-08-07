package service

import (
	"context"
	"errors"
	"mqtt-golang-rainfall-prediction/models"
	"mqtt-golang-rainfall-prediction/repository"
	"sort"
	"time"

	"github.com/mattn/go-tflite"
)

type Service interface {
	InsertData(ctx context.Context, data models.SensorData) error
	GetData(ctx context.Context) ([]models.SensorData, error)
	GetDataByDay(ctx context.Context) ([]models.SensorDataByDay, error)
	ProsesMessage(ctx context.Context, data models.SystemInfo) (models.SystemInfo, error)
	GetModelResult(ctx context.Context, input []float32) (float32, error)
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
		t1, err1 := time.Parse("01/02/2006", data[i].FormattedTime)
		t2, err2 := time.Parse("01/02/2006", data[j].FormattedTime)
		if err1 != nil || err2 != nil {
			return false // Handle parse error, maybe log or handle differently
		}
		return t1.After(t2) // Use After to sort descending
	})

	return data, nil
}

func (s *service) GetModelResult(ctx context.Context, input []float32) (float32, error) {
	model := tflite.NewModelFromFile("../model_lstm.tflite")
	if model == nil {
		return 0, errors.New("failed to load model")
	}
	defer model.Delete()

	options := tflite.NewInterpreterOptions()
	interpretter := tflite.NewInterpreter(model, options)

	if interpretter == nil {
		return 0, errors.New("failed to create interpreter")
	}
	defer interpretter.Delete()

	interpretter.AllocateTensors()

	inputTensor := interpretter.GetInputTensor(0)

	inputTensor.CopyFromBuffer(input)

	interpretter.Invoke()

	outputTensor := interpretter.GetOutputTensor(0)
	output := make([]float32, outputTensor.Dim(1))

	outputTensor.CopyToBuffer(&output[0])

	return output[0], nil

}
