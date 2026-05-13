package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/afidyoga/cinema-api/internal/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *model.User) error {
	q := `INSERT INTO users (id, name, email, password_hash, phone, role, created_at, updated_at)
          VALUES (:id, :name, :email, :password_hash, :phone, :role, :created_at, :updated_at)`
	_, err := r.db.NamedExec(q, u)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE email = $1`, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &u, err
}

type ScheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(s *model.Schedule) error {
	q := `INSERT INTO schedules (id, movie_id, studio_id, show_time, price, status, created_at, updated_at)
          VALUES (:id, :movie_id, :studio_id, :show_time, :price, :status, :created_at, :updated_at)`
	_, err := r.db.NamedExec(q, s)
	return err
}

func (r *ScheduleRepository) FindAll(page, limit int) ([]model.Schedule, int, error) {
	offset := (page - 1) * limit
	var total int
	if err := r.db.Get(&total, `SELECT COUNT(*) FROM schedules`); err != nil {
		return nil, 0, err
	}

	q := `SELECT s.*, m.title AS movie_title, st.name AS studio_name, c.name AS cinema_name
          FROM schedules s
          JOIN movies m ON m.id = s.movie_id
          JOIN studios st ON st.id = s.studio_id
          JOIN cinemas c ON c.id = st.cinema_id
          ORDER BY s.show_time ASC
          LIMIT $1 OFFSET $2`

	var schedules []model.Schedule
	err := r.db.Select(&schedules, q, limit, offset)
	return schedules, total, err
}

func (r *ScheduleRepository) FindByID(id string) (*model.Schedule, error) {
	var s model.Schedule
	q := `SELECT s.*, m.title AS movie_title, st.name AS studio_name, c.name AS cinema_name
          FROM schedules s
          JOIN movies m ON m.id = s.movie_id
          JOIN studios st ON st.id = s.studio_id
          JOIN cinemas c ON c.id = st.cinema_id
          WHERE s.id = $1`
	err := r.db.Get(&s, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &s, err
}

func (r *ScheduleRepository) Update(id string, req *model.UpdateScheduleRequest) (*model.Schedule, error) {
	existing, err := r.FindByID(id)
	if err != nil || existing == nil {
		return nil, err
	}

	if req.ShowTime != nil {
		existing.ShowTime = *req.ShowTime
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.UpdatedAt = time.Now()

	q := `UPDATE schedules SET show_time = :show_time, price = :price, status = :status, updated_at = :updated_at WHERE id = :id`
	_, err = r.db.NamedExec(q, existing)
	if err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *ScheduleRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM schedules WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
