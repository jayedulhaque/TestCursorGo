package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines user persistence (implemented implicitly by PostgresRepo).
type UserRepository interface {
	Create(ctx context.Context, name string) (*User, error)
	List(ctx context.Context) ([]User, error)
}

// EmployeeRepository defines employee persistence (implemented implicitly by PostgresRepo).
type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, name, email string) (*Employee, error)
	ListEmployees(ctx context.Context) ([]Employee, error)
}

// PostgresRepo implements UserRepository and EmployeeRepository using PostgreSQL via pgxpool.
type PostgresRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresRepo returns a repository backed by the given pool.
func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

// Create inserts a user and returns the stored row.
func (r *PostgresRepo) Create(ctx context.Context, name string) (*User, error) {
	const q = `
		INSERT INTO users (name)
		VALUES ($1)
		RETURNING id, name, created_at
	`
	var u User
	if err := r.pool.QueryRow(ctx, q, name).Scan(&u.ID, &u.Name, &u.CreatedAt); err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}
	return &u, nil
}

// List returns all users ordered by id.
func (r *PostgresRepo) List(ctx context.Context) ([]User, error) {
	const q = `SELECT id, name, created_at FROM users ORDER BY id`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return out, nil
}

// CreateEmployee inserts an employee and returns the stored row.
func (r *PostgresRepo) CreateEmployee(ctx context.Context, name, email string) (*Employee, error) {
	const q = `
		INSERT INTO employees (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at
	`
	var e Employee
	if err := r.pool.QueryRow(ctx, q, name, email).Scan(&e.ID, &e.Name, &e.Email, &e.CreatedAt); err != nil {
		return nil, fmt.Errorf("insert employee: %w", err)
	}
	return &e, nil
}

// ListEmployees returns all employees ordered by id.
func (r *PostgresRepo) ListEmployees(ctx context.Context) ([]Employee, error) {
	const q = `SELECT id, name, email, created_at FROM employees ORDER BY id`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query employees: %w", err)
	}
	defer rows.Close()

	var out []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan employee: %w", err)
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate employees: %w", err)
	}
	return out, nil
}
