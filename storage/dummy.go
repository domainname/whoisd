package storage

import (
	"errors"
)

// DummyRecord - standard record (struct) for dummy storage package
type DummyRecord struct {
	Host      string
	Port      int
	IndexBase string
	TypeTable string
}

// Search data in the storage
func (dummy *DummyRecord) Search(name string, query string) (map[string][]string, error) {

	result, err := dummy.searchRaw(dummy.TypeTable, name, query)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}

	data := make(map[string][]string)
	return data, nil
}

// SearchRelated - search data in the storage from related type or table
func (dummy *DummyRecord) SearchRelated(typeTable string, name string, query string) (map[string][]string, error) {

	result, err := dummy.searchRaw(typeTable, name, query)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}

	data := make(map[string][]string)
	return data, nil
}

// SearchMultiple - search multiple records of data in the storage
func (dummy *DummyRecord) SearchMultiple(typeTable string, name string, query string) (map[string][]string, error) {

	result, err := dummy.searchRaw(typeTable, name, query)
	if err != nil {
		return nil, err
	}

	data := make(map[string][]string)

	if len(result) > 0 {
		for _, item := range result {
			for key, value := range item {
				data[key] = append(data[key], value...)
			}
		}
		return data, nil
	}

	return data, nil
}

// search raw data in the storage
func (dummy *DummyRecord) searchRaw(typeTable string, name string, query string) ([]map[string][]string, error) {

	if len(typeTable) == 0 || len(name) == 0 || len(query) == 0 {
		return nil, errors.New("Incomplete request, request parameters could not be empty")
	}

	var data []map[string][]string

	for _, result := range dummyData[typeTable] {
		for _, item := range result[name] {
			if item == query {
				data = append(data, result)
				break
			}
		}
	}

	return data, nil
}
