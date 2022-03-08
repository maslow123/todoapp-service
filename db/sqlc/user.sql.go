// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    name,
    address,
    pic,
    hashed_password,
    email
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, name, address, pic, created_at, updated_at, hashed_password, email
`

type CreateUserParams struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	Pic            string `json:"pic"`
	HashedPassword string `json:"hashed_password"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Name,
		arg.Address,
		arg.Pic,
		arg.HashedPassword,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Pic,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.Email,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, name, address, pic, created_at, updated_at, hashed_password, email FROM users
WHERE email = $1
`

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Pic,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.Email,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, address, pic, created_at, updated_at, hashed_password, email FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.Pic,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.HashedPassword,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET name = $2, address = $3, pic = $4, email = $5, updated_at = now()
WHERE id = $1
RETURNING id, name, address, pic, created_at, updated_at, hashed_password, email
`

type UpdateUserParams struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Pic     string `json:"pic"`
	Email   string `json:"email"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Address,
		arg.Pic,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Pic,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.Email,
	)
	return i, err
}

const updateUserPhoto = `-- name: UpdateUserPhoto :one
UPDATE users
SET pic = $2, updated_at = now()
WHERE email = $1
RETURNING id, name, address, pic, created_at, updated_at, hashed_password, email
`

type UpdateUserPhotoParams struct {
	Email string `json:"email"`
	Pic   string `json:"pic"`
}

func (q *Queries) UpdateUserPhoto(ctx context.Context, arg UpdateUserPhotoParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPhoto, arg.Email, arg.Pic)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.Pic,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedPassword,
		&i.Email,
	)
	return i, err
}
