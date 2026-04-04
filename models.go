package main

import "time"

// User is the persisted user model exposed via JSON.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest is the body for POST /users.
type CreateUserRequest struct {
	Name string `json:"name"`
}

// Employee is the persisted employee model exposed via JSON.
type Employee struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateEmployeeRequest is the body for POST /employees.
type CreateEmployeeRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
