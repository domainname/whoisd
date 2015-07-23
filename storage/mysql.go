package storage

import (
	"errors"
)

// MysqlRecord - standard record (struct) for mysql storage package
type MysqlRecord struct {
	Host     string
	Port     int
	DataBase string
	Table    string
}

// Search data in the storage
// TODO - Mysql storage is not released
func (mysql *MysqlRecord) Search(name string, query string) (map[string][]string, error) {

	return nil, errors.New("Mysql driver not released")
}

// SearchRelated - search data in the storage from related type or table
// TODO - Mysql storage is not released
func (mysql *MysqlRecord) SearchRelated(typeTable string, name string, query string) (map[string][]string, error) {

	return nil, errors.New("Mysql driver not released")
}

// SearchMultiple - search multiple records of data in the storage
// TODO - Mysql storage is not released
func (mysql *MysqlRecord) SearchMultiple(typeTable string, name string, query string) (map[string][]string, error) {

	return nil, errors.New("Mysql driver not released")
}
