package postgres

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/lib/pq"
)

const (
	CodeUniqueViolation              = "23505"
	CodeNotNullViolation             = "23502"
	CodeIntegrityConstraintViolation = "23000"
)

var (
	colPat = regexp.MustCompile(`Key \((.+)\)=`)
	valPat = regexp.MustCompile(`Key \(.+\)=\((.+)\)`)
)

type DBError struct {
	Message    string
	Code       string
	Detail     string
	Constraint string
}

type Constraint struct {
	Name string
}

type Constraints struct {
	Map map[string]*Constraint
	Mux sync.RWMutex
}

func (c *Constraints) Add(name string) error {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	if _, present := c.Map[name]; present {
		return fmt.Errorf("Constaint %s already added", name)
	}
	c.Map[name] = &Constraint{Name: name}
	return nil
}

func (dbe *DBError) Error() string {
	return dbe.Message
}

func findValue(s string) string {
	fmt.Println(s)
	r := valPat.FindStringSubmatch(s)
	if len(r) > 0 {
		return r[1]
	}
	return ""
}

func findCol(s string) string {
	r := colPat.FindStringSubmatch(s)
	if len(r) > 0 {
		return r[1]
	}
	return ""
}

func ParseError(err error) error {
	if err == nil {
		return nil
	}

	switch pqErr := err.(type) {
	case *pq.Error:
		switch pqErr.Code {
		case CodeUniqueViolation:
			var msg string
			col := findCol(pqErr.Detail)
			val := findValue(pqErr.Detail)
			if len(col) == 0 {
				msg = fmt.Sprintf("Col already contain value (%s)", val)
			} else {
				msg = fmt.Sprintf("Col (%s) already contain value (%s)", col, val)
			}
			return &DBError{
				Message:    msg,
				Code:       string(pqErr.Code),
				Detail:     pqErr.Detail,
				Constraint: pqErr.Constraint,
			}

		case CodeIntegrityConstraintViolation:
			panic("Not Implemented")

		case CodeNotNullViolation:
			msg := fmt.Sprintf("Column (%s) cannot be left blank", pqErr.Column)
			return &DBError{
				Message:    msg,
				Code:       string(pqErr.Code),
				Detail:     pqErr.Detail,
				Constraint: pqErr.Constraint,
			}
		}

	default:
		return pqErr
	}
	return nil
}
