package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"orchestrator/internal/domain/models"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf(
			"DATA LAYER: storage.postgres.New: couldn't open a database: %w",
			err,
		)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) SaveOperation(
	ctx context.Context,
	settings models.Operation,
) error {
	return nil
}

func (s *Storage) GetOperation(
	ctx context.Context,
	operation string,
) (models.Operation, error) {
	return models.Operation{
		Id:           "1",
		Operation:    "(2+2)*2",
		CreationAt:   "2024-02-07",
		CalculatedAt: "2024-02-07",
	}, nil
}

func (s *Storage) SaveOperationExecutionTime(
	ctx context.Context,
	settings models.Settings,
) error {
	return nil
}

func (s *Storage) GetOperationExecutionTime(
	ctx context.Context,
	value any,
) (models.Settings, error) {
	return models.Settings{
		Id:                          "1",
		PlusOperationExecutionTime:  20,
		MinusOperationExecutionTime: 202,
		MultiplicationExecutionTime: 302,
		DivisionExecutionTime:       400,
	}, nil
}

//// SaveUser saves user to db.
//func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
//	var id int
//	query := "INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id"
//	err := s.db.QueryRowContext(ctx, query, email, passHash).Scan(&id)
//	if err != nil {
//		return 0, fmt.Errorf(
//			"DATA LAYER: storage.postgres.SaveUser: couldn't save user  %w",
//			err,
//		)
//	}
//	return int64(id), nil
//}
//
//func (s *Storage) GetUser(ctx context.Context, value any) (models.User, error) {
//	var row *sql.Row
//	switch sqlParam := value.(type) {
//	case int:
//		query := "SELECT id, email, pass_hash, is_admin FROM users WHERE (id = $1);"
//		row = s.db.QueryRowContext(ctx, query, sqlParam)
//	case string:
//		query := "SELECT id, email, pass_hash, is_admin FROM users WHERE (email = $1);"
//		row = s.db.QueryRowContext(ctx, query, sqlParam)
//	default:
//		return models.User{}, fmt.Errorf(
//			"DATA LAYER: storage.postgres.GetUser: %w",
//			storage.ErrWrongParamType,
//		)
//	}
//
//	var user models.User
//	err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.IsAdmin)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return models.User{}, fmt.Errorf(
//				"DATA LAYER: storage.postgres.GetUser: %w",
//				storage.ErrUserNotFound,
//			)
//		}
//		return models.User{}, fmt.Errorf(
//			"DATA LAYER: storage.postgres.GetUser: %w",
//			err,
//		)
//	}
//	return user, nil
//}
