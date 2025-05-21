package database

import (
	"github.com/pgvector/pgvector-go"
)

func ToPgVector(v []float64) pgvector.Vector {
	return pgvector.NewVector(convertFloat64ToFloat32(v))
}

func convertFloat64ToFloat32(input []float64) []float32 {
	if input == nil {
		return nil
	}
	output := make([]float32, len(input))
	for i, val64 := range input {
		output[i] = float32(val64) // Individual element conversion
	}
	return output
}
