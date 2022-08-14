package sdkcm

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	dateFmt                       = "2006-01-02"
	timeAddDuration time.Duration = 0
)

func SetTimeZone(zone int64) {
	timeAddDuration = time.Hour * time.Duration(zone)
}

type JSONDate time.Time

// Implement method MarshalJSON to output date with in formatted
func (d JSONDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(d).Add(timeAddDuration).Format(dateFmt))
	return []byte(stamp), nil
}

func (d *JSONDate) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(dateFmt, strings.Replace(string(data), "\"", "", -1))

	*d = JSONDate(t)

	if err != nil {
		return err
	}
	return nil
}

// This method for mapping JSONDate to date data type in sql
func (d *JSONDate) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return time.Time(*d).Format(dateFmt), nil
}

// This method for scanning JSONDate from date data type in sql
func (d *JSONDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if v, ok := value.(time.Time); ok {
		*d = JSONDate(v)
		return nil
	}

	return errors.New("invalid Scan Source")
}

func (d JSONDate) GetBSON() (interface{}, error) {
	return time.Time(d), nil
}
