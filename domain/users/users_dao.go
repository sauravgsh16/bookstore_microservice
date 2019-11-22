// DAO - domain access object: Provides the means to access the persistance layers
// No other layer in the application is responsible for accessing the db

package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/sauravgsh16/bookstore_users-api/datasource/postgres/usersdb"
	"github.com/sauravgsh16/bookstore_users-api/utils/errors"
	"github.com/sauravgsh16/bookstore_users-api/utils/postgres"
)

const (
	queryInsertUser       = `INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES($1, $2, $3, $4, $5, $6) RETURNING ID;`
	querySelectUser       = `SELECT ID, first_name, last_name, email, date_created, status FROM users WHERE ID=($1);`
	queryUpdateuser       = `UPDATE users SET first_name=($1), last_name=($2), email=($3) WHERE ID=($4);`
	queryDeleteUser       = `DELETE FROM users WHERE ID=($1);`
	queryFindUserByStatus = `SELECT ID, FIRST_NAME, LAST_NAME, EMAIL, DATE_CREATED, STATUS FROM users WHERE STATUS=($1);`
)

var (
	userDB = make(map[int64]*User)
)

func getConn() (*sql.Conn, context.Context) {
	ctx := context.Background()
	return usersdb.DB.GetConn(ctx), ctx
}

func handleDBError(err error) *errors.RestErr {
	err = postgres.ParseError(err)
	dbErr, ok := err.(*pq.Error)
	if ok {
		return errors.NewBadRequestError(dbErr.Message)
	}
	return nil
}

// Get populates the user pointer or returns error if not found
func (u *User) Get(userID int64) *errors.RestErr {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, querySelectUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	row := stmt.QueryRowContext(ctx, userID)

	if err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
		if err := handleDBError(err); err != nil {
			return err
		}
		return errors.NewNotFoundError(fmt.Sprintf("failed to retrieve rows: %s", err.Error()))
	}
	return nil
}

// Save the user to the db
func (u *User) Save() *errors.RestErr {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	var returnedID int64

	if err = stmt.QueryRowContext(ctx, u.FirstName, u.LastName, u.Email, u.DateCreated, u.Status, u.Password).Scan(&returnedID); err != nil {
		if err := handleDBError(err); err != nil {
			return err
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	u.ID = returnedID
	return nil
}

// Update the user in the db
func (u *User) Update() *errors.RestErr {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, queryUpdateuser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	if _, err = stmt.ExecContext(ctx, u.FirstName, u.LastName, u.Email, u.ID); err != nil {
		if err := handleDBError(err); err != nil {
			return err
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to update user: %s", err.Error()))
	}
	return nil
}

// Delete user from db
func (u *User) Delete() *errors.RestErr {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	if _, err = stmt.ExecContext(ctx, u.ID); err != nil {
		if err := handleDBError(err); err != nil {
			return err
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to delete user: %s", err.Error()))
	}
	return nil
}

// FindByStatus retusn a list of user where status is passed as an agrument
func (u *User) FindByStatus(status string) ([]*User, *errors.RestErr) {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	rows, err := stmt.QueryContext(ctx, status)
	if err != nil {
		if err := handleDBError(err); err != nil {
			return nil, err
		}
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		u := new(User)
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
			if err := handleDBError(err); err != nil {
				return nil, err
			}
			return nil, errors.NewInternalServerError(err.Error())
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	if len(users) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("Users with status - %s - not found", status))
	}
	return users, nil
}
