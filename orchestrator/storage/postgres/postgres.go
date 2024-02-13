package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"orchestrator/internal/domain/models"
	"orchestrator/storage"
	"time"
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
	query := "INSERT INTO operations(uid, operation, result, status, created_at) VALUES($1, $2, $3, $4, $5)"
	_, err := s.db.ExecContext(ctx, query, operation.Id, operation.Operation, value, "process", time.Now())
	if err != nil {
		return fmt.Errorf(
			"DATA LAYER: storage.postgres.SaveOperation: couldn't save Operation  %w",
			err,
		)
	}
	return nil
}

func (s *Storage) UpdateOperation(
	ctx context.Context,
	operation models.Operation,
) error {
	query := "UPDATE operations SET result = $1, status = $2, calculated_at = $3 WHERE uid = $4;"
	_, err := s.db.ExecContext(ctx, query, operation.Result, operation.Status, time.Now(), operation.Id)
	if err != nil {
		return fmt.Errorf(
			"DATA LAYER: storage.postgres.UpdateOperation: couldn't update Operation  %w",
			err,
		)
	}
	return nil
}

func (s *Storage) GetOperation(
	ctx context.Context,
	operation string,
) (models.Operation, error) {

	query := "SELECT uid, operation, result, created_at, calculated_at FROM operations WHERE (operation = $1);"
	row := s.db.QueryRowContext(ctx, query, operation)

	var foundOperation models.Operation
	err := row.Scan(&foundOperation.Id, &foundOperation.Operation, &foundOperation.Result, &foundOperation.CreatedAt, &foundOperation.CalculatedAt)
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

func (s *Storage) GetOperationById(
	ctx context.Context,
	uid string,
) (models.Operation, error) {

	query := "SELECT uid, operation, result, status,  created_at, calculated_at FROM operations WHERE (uid = $1);"
	row := s.db.QueryRowContext(ctx, query, uid)

	var foundOperation models.Operation
	err := row.Scan(&foundOperation.Id, &foundOperation.Operation, &foundOperation.Result, &foundOperation.Status, &foundOperation.CreatedAt, &foundOperation.CalculatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return foundOperation, fmt.Errorf(
				"DATA LAYER: storage.postgres.GetOperationById: %w",
				storage.ErrOperationNotFound,
			)
		}
		return foundOperation, fmt.Errorf(
			"DATA LAYER: storage.postgres.GetOperationById: %w",
			err,
		)
	}
	return foundOperation, nil
}

func (s *Storage) UpdateSettingsExecutionTime(
	ctx context.Context,
	opType storage.OperationType,
	executionTime int,
) error {

	var fieldName string
	switch opType {
	case storage.PlusOperation:
		fieldName = "plus_operation_execution_time"
	case storage.MinusOperation:
		fieldName = "minus_operation_execution_time"
	case storage.MultiplicationOperation:
		fieldName = "multiplication_operation_execution_time"
	case storage.DivisionOperation:
		fieldName = "division_operation_execution_time"
	default:
		//TODO: use storage errors
		return errors.New("Unknown operation type")
	}
	fmt.Println("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	query := fmt.Sprintf("UPDATE settings SET %s = $1 WHERE id = 1;", fieldName)
	_, err := s.db.ExecContext(ctx, query, executionTime)
	if err != nil {
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", err)
		return fmt.Errorf("DATA LAYER: storage.postgres.UpdateSettingsExecutionTime: couldn't update %s operation execution time %w", fieldName, err)
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
