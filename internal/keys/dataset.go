package keys

import (
	"crypto/rand"
)

func GenerateDatasetKey() ([]byte, error) {
	datasetKey := make([]byte, 32)
	_, err := rand.Read(datasetKey)
	if err != nil {
		return nil, err
	}
	return datasetKey, nil
}
