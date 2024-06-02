package user

import (
	"context"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/user/entities"
	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/pkg/logger"
	"github.com/tamboto2000/99-backend-exercise/pkg/sqli"
)

type UserRepository interface {
	Create(ctx context.Context, user entities.User) error
	GetByID(ctx context.Context, id int64) (entities.User, error)
	GetAll(ctx context.Context, offset, limit int) ([]entities.User, error)
}

type userRepo struct {
	db *sqli.DB
}

func NewUserRepository(db *sqli.DB) UserRepository {
	return &userRepo{db: db}
}

func (u *userRepo) Create(ctx context.Context, user entities.User) error {
	q := `
	INSERT INTO users (
		id,
		name,
		created_at,
		updated_at
	) 
	VALUES (
		$1,
		$2,
		$3,
		$4
	)`

	_, err := u.db.Exec(ctx, q, user.ID, user.Name, user.CreatedAt, user.UpdatedAt)

	return err
}

func (u *userRepo) GetByID(ctx context.Context, id int64) (entities.User, error) {
	var user entities.User
	q := `
	SELECT
		id,
		name,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`
	row := u.db.QueryRow(ctx, q, id)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sqli.ErrNoRows {
			return user, errors.New("user not found", errors.CodeNotExists)
		}

		logger.Error(err.Error())
		return user, err
	}

	return user, nil
}

func (u *userRepo) GetAll(ctx context.Context, offset, limit int) ([]entities.User, error) {
	q := `
	SELECT
		id,
		name,
		created_at,
		updated_at
	FROM users
	ORDER BY created_at DESC
	OFFSET $1
	LIMIT $2
	`

	var users []entities.User
	rows, err := u.db.Query(ctx, q, offset, limit)
	if err != nil {
		if err == sqli.ErrNoRows {
			return nil, nil
		}

		logger.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		var user entities.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			if err == sqli.ErrNoRows {
				return nil, nil
			}

			logger.Error(err.Error())
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
