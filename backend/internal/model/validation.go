package model

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

func IsDateValid() validation.RuleFunc {
	return func(value interface{}) error {
		tmp, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("incorrect validating field type")
		}

		if tmp.IsZero() {
			return fmt.Errorf("date can't be default")
		}

		if !tmp.Before(time.Now()) {
			return fmt.Errorf("date can't be in future")
		}
		return nil
	}
}