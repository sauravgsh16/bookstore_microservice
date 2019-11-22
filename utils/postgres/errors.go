package postgres

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/lib/pq"
)

const (
	codeUniqueViolation              = "23505"
	codeNotNullViolation             = "23502"
	codeIntegrityConstraintViolation = "23000"
)

var (
	colPat = regexp.MustCompile(`Key \((.+)\)=`)
	valPat = regexp.MustCompile(`Key \(.+\)=\((.+)\)`)
)

// DBError struct
type DBError struct {
	Message    string
	Code       string
	Detail     string
	Constraint string
}

// Constraint struct
type Constraint struct {
	Name string
}

// Constraints contains map of user defined constraints
type Constraints struct {
	Map map[string]*Constraint
	Mux sync.RWMutex
}

// Add new constraint
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

// ParseError parses db errors
func ParseError(err error) error {
	if err == nil {
		return nil
	}

	switch pqErr := err.(type) {
	case *pq.Error:
		switch pqErr.Code {
		case codeUniqueViolation:
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

		case codeIntegrityConstraintViolation:
			panic("Not Implemented")

		case codeNotNullViolation:
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
