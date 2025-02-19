package csv

import (
	"encoding/csv"
	"fmt"
	"io"
)

type MapperFunc[T any] func([]string) (T, error)

func Decode[T any](reader io.Reader, mapper MapperFunc[T]) ([]T, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV data: %w", err)
	}

	var data []T

	for i, record := range records {
		if i == 0 {
			continue
		}
		item, err := mapper(record)
		if err != nil {
			return nil, fmt.Errorf("error mapping record: %w", err)
		}
		data = append(data, item)
	}

	return data, nil
}
