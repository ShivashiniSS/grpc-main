package studentdex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("student not found")
var ErrAlreadyExists = errors.New("this student already exists")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Student) MarshalJSON() ([]byte, error) {
	type Alias Student
	return json.Marshal(&struct {
		*Alias
		CreatedAtUnix int64 `json:"created_at_unix"`
		UpdatedAtUnix int64 `json:"updated_at_unix"`
	}{
		Alias:         (*Alias)(s),
		CreatedAtUnix: s.CreatedAt.Unix(),
		UpdatedAtUnix: s.UpdatedAt.Unix(),
	})
}

func (s *Student) UnmarshalJSON(data []byte) error {
	type Alias Student
	aux := &struct {
		*Alias
		CreatedAtUnix int64 `json:"created_at_unix"`
		UpdatedAtUnix int64 `json:"updated_at_unix"`
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.CreatedAt = time.Unix(aux.CreatedAtUnix, 0)
	s.UpdatedAt = time.Unix(aux.UpdatedAtUnix, 0)
	return nil
}

func (p *Repository) Insert(ctx context.Context, studentDetails Student) error {
	_, err := p.db.Exec(
		ctx,
		"INSERT INTO student (id, name, grade, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		studentDetails.ID,
		studentDetails.Name,
		studentDetails.Grade,
		studentDetails.CreatedAt,
		studentDetails.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert student: %w", err)
	}

	return nil
}

func (p *Repository) FindAll(ctx context.Context) ([]Student, error) {
	rows, err := p.db.Query(ctx, "SELECT id, name, grade, created_at, updated_at FROM student")
	if err != nil {
		return nil, fmt.Errorf("failed to find all student details: %w", err)
	}

	defer rows.Close()
	res := make([]Student, 0)

	for rows.Next() {
		var studentDetails Student

		err := rows.Scan(
			&studentDetails.ID,
			&studentDetails.Name,
			&studentDetails.Grade,
			&studentDetails.CreatedAt,
			&studentDetails.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to find student: %w", err)
		}

		res = append(res, studentDetails)
	}

	return res, nil
}

func (p *Repository) FindByID(ctx context.Context, id int32) (Student, error) {
	res := Student{}

	err := p.db.QueryRow(
		ctx,
		"SELECT id, name, grade, created_at, updated_at FROM student WHERE id = $1",
		id,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Grade,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Student{}, ErrNotFound
	} else if err != nil {
		return Student{}, fmt.Errorf("failed to find student by id: %w", err)
	}

	return res, nil
}

func (p *Repository) Update(ctx context.Context, studentDetails Student) error {
	res, err := p.db.Exec(
		ctx,
		"UPDATE student SET name = $1, grade = $2, updated_at = $3 WHERE id = $4",
		studentDetails.Name,
		studentDetails.Grade,
		studentDetails.UpdatedAt,
		studentDetails.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
