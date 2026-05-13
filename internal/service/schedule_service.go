package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/afidyoga/cinema-api/internal/model"
	"github.com/afidyoga/cinema-api/internal/repository"
)

type ScheduleService struct {
	repo *repository.ScheduleRepository
}

func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) Create(req *model.CreateScheduleRequest) (*model.Schedule, error) {
	now := time.Now()
	schedule := &model.Schedule{
		ID:        uuid.NewString(),
		MovieID:   req.MovieID,
		StudioID:  req.StudioID,
		ShowTime:  req.ShowTime,
		Price:     req.Price,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(schedule); err != nil {
		return nil, err
	}
	return s.repo.FindByID(schedule.ID)
}

func (s *ScheduleService) GetAll(page, limit int) ([]model.Schedule, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.FindAll(page, limit)
}

func (s *ScheduleService) GetByID(id string) (*model.Schedule, error) {
	sch, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if sch == nil {
		return nil, errors.New("schedule not found")
	}
	return sch, nil
}

func (s *ScheduleService) Update(id string, req *model.UpdateScheduleRequest) (*model.Schedule, error) {
	updated, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errors.New("schedule not found")
	}
	return updated, nil
}

func (s *ScheduleService) Delete(id string) error {
	err := s.repo.Delete(id)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("schedule not found")
	}
	return err
}
