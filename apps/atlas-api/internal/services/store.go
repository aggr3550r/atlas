package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateServiceParams struct {
	Name    string
	Version string
	Status  string
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, params CreateServiceParams) (Service, error) {
	service := Service{
		ID:      uuid.New(),
		Name:    params.Name,
		Version: params.Version,
		Status:  params.Status,
	}

	err := s.db.QueryRow(ctx, `
		INSERT INTO services (id, name, version, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, version, status, created_at, updated_at
	`, service.ID, service.Name, service.Version, service.Status).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.Status,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		return Service{}, err
	}

	return service, nil
}

func (s *Store) List(ctx context.Context) ([]Service, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, name, version, status, created_at, updated_at
		FROM services
		ORDER BY created_at ASC, name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var service Service
		if err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Version,
			&service.Status,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (s *Store) Get(ctx context.Context, id uuid.UUID) (Service, error) {
	var service Service
	err := s.db.QueryRow(ctx, `
		SELECT id, name, version, status, created_at, updated_at
		FROM services
		WHERE id = $1
	`, id).Scan(
		&service.ID,
		&service.Name,
		&service.Version,
		&service.Status,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		return Service{}, err
	}

	return service, nil
}
