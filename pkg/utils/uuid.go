//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package utils

import "github.com/google/uuid"

type UUID struct{}

func NewUUID() *UUID {
	return &UUID{}
}

type UUIDInterface interface {
	Get() (string, error)
}

// インターフェースを満たしているかを確認
var _ UUIDInterface = (*UUID)(nil)

func (*UUID) Get() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
