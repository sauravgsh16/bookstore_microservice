// DAO - domain access object: Provides the means to access the persistance layers
// No other layer in the application is responsible for accessing the db

package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sauravgsh16/bookstore_users-api/logger"

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
	queryFindByEmailPwd   = `SELECT ID, FIRST_NAME, LAST_NAME, EMAIL, DATE_CREATED, STATUS FROM users WHERE EMAIL=($1) AND PASSWORD=($2);`
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
		logger.Error("Query failed with error: ", err)
		return errors.NewBadRequestError(dbErr.Message)
	}
	return nil
}

// Get populates the user pointer or returns error if not found
func (u *User) Get(userID int) *errors.RestErr {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, querySelectUser)
	if err != nil {
		logger.Error("failed to prepare statement: ", err)
		return errors.NewInternalServerError("database error")
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
		logger.Error("failed to prepare statement: ", err)
		return errors.NewInternalServerError("database error")
	}

	var returnedID int

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
		logger.Error("failed to prepare statement: ", err)
		return errors.NewInternalServerError("database error")
	}

	if _, err = stmt.ExecContext(ctx, u.FirstName, u.LastName, u.Email, u.ID); err != nil {
		if err := handleDBError(err); err != nil {
			return err
		}
		logger.Error("failed to execute update query, error: ", err)
		return errors.NewInternalServerError("database error when trying to update")
	}
	return nil
}

// FindByEmailPassword finds user by email and password
func (u *User) FindByEmailPassword(email, pwd string) *errors.RestErr {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, queryFindByEmailPwd)
	if err != nil {
		logger.Error("failed to prepare statement: ", err)
	}

	fmt.Printf("%#v\n", stmt)

	row := stmt.QueryRowContext(ctx, email, pwd)

	fmt.Printf("%#v\n", row)

	if err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
		if err := handleDBError(err); err != nil {
			return err
		}
		return errors.NewNotFoundError(fmt.Sprintf("failed to retrieve rows: %s", err.Error()))
	}
	return nil
}

// Delete user from db
func (u *User) Delete() *errors.RestErr {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, queryDeleteUser)
	if err != nil {
		logger.Error("failed to prepare statement: ", err)
		return errors.NewInternalServerError("database error")
	}

	if _, err = stmt.ExecContext(ctx, u.ID); err != nil {
		if err := handleDBError(err); err != nil {
			return err
		}
		logger.Error("failed to execute delete query, error: ", err)
		return errors.NewInternalServerError("database error when trying to delete user")
	}
	return nil
}

// FindByStatus retusn a list of user where status is passed as an agrument
func (u *User) FindByStatus(status string) ([]*User, *errors.RestErr) {
	conn, ctx := getConn()
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, queryFindUserByStatus)
	if err != nil {
		logger.Error("failed to prepare statement: ", err)
		return nil, errors.NewInternalServerError("database error")
	}

	rows, err := stmt.QueryContext(ctx, status)
	if err != nil {
		if err := handleDBError(err); err != nil {
			return nil, err
		}
		logger.Error("failed to execute delete query, error: ", err)
		return nil, errors.NewInternalServerError("database error when trying to execute query")
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		u := new(User)
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
			if err := handleDBError(err); err != nil {
				return nil, err
			}
			logger.Error("failed to scan rows, error: ", err)
			return nil, errors.NewInternalServerError("database error")
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		logger.Error("Row error: ", err)
		return nil, errors.NewInternalServerError("database error")
	}
	if len(users) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("Users with status - %s - not found", status))
	}
	return users, nil
}
