package utility

import (
	"github.com/google/uuid"
)

func NewRandomUUID() string {
	return uuid.NewString()
}
