package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"orchestrator/internal/domain/models"
	"orchestrator/storage"
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
	operation models.Operation,
	value any,
) error {
	query := "INSERT INTO operations(uid, operation, result) VALUES($1, $2, $3)"
	_, err := s.db.ExecContext(ctx, query, operation.Id, operation.Operation, value)
	if err != nil {
		return fmt.Errorf(
			"DATA LAYER: storage.postgres.SaveOperation: couldn't save user  %w",
			err,
		)
	}
	return nil
}

func (s *Storage) GetOperation(
	ctx context.Context,
	operation string,
) (models.Operation, error) {

	query := "SELECT uid, operation, result, creation_at, calculated_at FROM operations WHERE (operation = $1);"
	row := s.db.QueryRowContext(ctx, query, operation)

	var foundOperation models.Operation
	err := row.Scan(&foundOperation.Id, &foundOperation.Operation, &foundOperation.Result, &foundOperation.CreationAt, &foundOperation.CalculatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return foundOperation, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetOperation: %w",
				storage.ErrOperationNotFound,
			)
		}
		return foundOperation, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetOperation: %w",
			err,
		)
	}
	return foundOperation, nil
}

func (s *Storage) SaveOperationExecutionTime(
	ctx context.Context,
	settings models.Settings,
) error {
	query := "INSERT INTO settings(id, plus_operation_execution_time, minus_operation_execution_time, multiplication_operation_execution_time, division_operation_execution_time ) VALUES($1, $2)"
	_, err := s.db.ExecContext(
		ctx,
		query,
		settings.Id,
		settings.PlusOperationExecutionTime,
		settings.MinusOperationExecutionTime,
		settings.MultiplicationExecutionTime,
		settings.DivisionExecutionTime,
	)
	if err != nil {
		return fmt.Errorf(
			"DATA LAYER: storage.postgres.SaveOperation: couldn't save user  %w",
			err,
		)
	}
	return nil
}

func (s *Storage) GetOperationExecutionTime(
	ctx context.Context,
) (models.Settings, error) {

	query := "SELECT id, plus_operation_execution_time, minus_operation_execution_time, multiplication_operation_execution_time, division_operation_execution_time FROM settings WHERE (id = 1);"
	row := s.db.QueryRowContext(ctx, query)

	var foundSettings models.Settings
	err := row.Scan(
		&foundSettings.Id,
		&foundSettings.PlusOperationExecutionTime,
		&foundSettings.MinusOperationExecutionTime,
		&foundSettings.MultiplicationExecutionTime,
		&foundSettings.DivisionExecutionTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return foundSettings, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetOperationExecutionTime: %w",
				storage.ErrOperationNotFound,
			)
		}
		return foundSettings, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetOperationExecutionTime: %w",
			err,
		)
	}
	return foundSettings, nil
}
