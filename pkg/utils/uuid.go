//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package utils

import "github.com/google/uuid"

type uuID struct{}

func NewUUID() UUID {
	return &uuID{}
}

type UUID interface {
	Get() (string, error)
}

// インターフェースを満たしているかを確認
var _ UUID = (*uuID)(nil)

func (*uuID) Get() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
