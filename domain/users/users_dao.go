// DAO - domain access object: Provides the means to access the persistance layers
// No other layer in the application is responsible for accessing the db

package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sauravgsh16/bookstore_users-api/datasource/postgres/usersdb"
	"github.com/sauravgsh16/bookstore_users-api/utils/dates"
	"github.com/sauravgsh16/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser = `INSERT INTO users(first_name, last_name, email, date_created) VALUES($1, $2, $3, $4) RETURNING ID;`
	querySelectUser = `SELECT ID, first_name, last_name, email, date_created FROM users WHERE ID=($1);`
)

var (
	userDB = make(map[int64]*User)
)

func getConn() (*sql.Conn, context.Context) {
	ctx := context.Background()
	return usersdb.DB.GetConn(ctx), ctx
}

// Save the user to the db
func (u *User) Save() *errors.RestErr {
	conn, ctx := getConn()

	stmt, err := conn.PrepareContext(ctx, queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer conn.Close()

	u.DateCreated = dates.GetNowString()

	var returnedID int64

	err = stmt.QueryRowContext(ctx, u.FirstName, u.LastName, u.Email, u.DateCreated).Scan(&returnedID)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return errors.NewBadRequestError(fmt.Sprintf("user \"%s\" already exists", u.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	u.ID = returnedID
	return nil
}

// Get populates the user pointer or returns error if not found
func (u *User) Get(userID int64) *errors.RestErr {
	conn, ctx := getConn()
	stmt, err := conn.PrepareContext(ctx, querySelectUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("failed to execute query: %s", err.Error()))
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return errors.NewInternalServerError(fmt.Sprintf("failed to retrieve rows from db: %s", err.Error()))
		}
		return errors.NewNotFoundError(fmt.Sprintf("User with ID: %d - not found", userID))
	}
	if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated); err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("failed to retrieve rows: %s", err.Error()))
	}
	return nil
}
