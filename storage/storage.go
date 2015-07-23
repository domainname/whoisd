package storage

import (
	"log"
	"sort"
	"strings"
	"time"

	"github.com/takama/whoisd/config"
	"github.com/takama/whoisd/mapper"
)

// Storage - the interface for every implementation of storage
type Storage interface {
	Search(name string, query string) (map[string][]string, error)
	SearchRelated(typeTable string, name string, query string) (map[string][]string, error)
	SearchMultiple(typeTable string, name string, query string) (map[string][]string, error)
}

// Record - standard record (struct) for storage package
type Record struct {
	CurrentStorage Storage
	Mapper         *mapper.Record
}

// New - returns new Storage instance
func New(conf *config.Record, mapp *mapper.Record) *Record {
	switch strings.ToLower(conf.Storage.StorageType) {
	case "mysql":
		return &Record{
			&MysqlRecord{
				conf.Storage.Host,
				conf.Storage.Port,
				conf.Storage.IndexBase,
				conf.Storage.TypeTable,
			},
			mapp,
		}
	case "elasticsearch":
		return &Record{
			&ElasticsearchRecord{
				conf.Storage.Host,
				conf.Storage.Port,
				conf.Storage.IndexBase,
				conf.Storage.TypeTable,
			},
			mapp,
		}
	case "dummy":
		fallthrough
	default:
		return &Record{
			&DummyRecord{
				conf.Storage.Host,
				conf.Storage.Port,
				conf.Storage.IndexBase,
				conf.Storage.TypeTable,
			},
			mapp,
		}
	}
}

// Search and sort a data from the storage
func (storage *Record) Search(query string) (answer string, ok bool) {
	ok = false
	answer = "not found\n"
	if len(strings.TrimSpace(query)) == 0 {
		log.Println("Empty query")
	} else {
		mapp, err := storage.LoadMapper(strings.TrimSpace(query))
		if err != nil {
			log.Println("Query:", query, err.Error())
		} else {
			if mapp == nil || mapp.Count() == 0 {
				return answer, ok
			}
			ok = true

			// get keys of a map and sort their
			keys := make([]string, 0, mapp.Count())
			for key := range mapp.Fields {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			answer = prepareAnswer(mapp, keys)
		}
	}

	return answer, ok
}

// LoadMapper - Loads a data into the Mapper
func (storage *Record) LoadMapper(query string) (*mapper.Record, error) {

	var err error

	mapp := new(mapper.Record)
	mapp.Fields = make(map[string]mapper.Field)
	baseRecord := make(map[string][]string)
	relatedRecord := make(map[string]map[string][]string)

	// Loads prearranged values
	for index, record := range storage.Mapper.Fields {
		if len(record.Value) != 0 && len(record.Related) == 0 &&
			len(record.RelatedBy) == 0 && len(record.RelatedTo) == 0 {
			mapp.Fields[index] = mapper.Field{
				Key:      record.Key,
				Value:    record.Value,
				Format:   record.Format,
				Multiple: record.Multiple,
				Hide:     record.Hide,
			}
		}
	}

	// Loads base record
	for index, record := range storage.Mapper.Fields {
		// Check for base record
		if len(record.Value) == 0 && len(record.Related) != 0 &&
			(len(record.RelatedBy) == 0 || len(record.RelatedTo) == 0) {
			// if not cached, do it
			if len(baseRecord) == 0 {
				baseRecord, err = storage.CurrentStorage.Search(record.Related, query)
				if err != nil {
					return nil, err
				}
				if len(baseRecord) == 0 {
					return nil, nil
				}
			}
			answer := []string{}

			// collects all fields into answer
			for _, item := range record.Name {
				if result, ok := baseRecord[item]; ok {
					answer = append(answer, result...)
				}
			}

			mapp.Fields[index] = mapper.Field{
				Key:      record.Key,
				Value:    answer,
				Format:   record.Format,
				Multiple: record.Multiple,
				Hide:     record.Hide,
			}
		}
	}

	// Loads related records
	for index, record := range storage.Mapper.Fields {
		// Check for related record
		if len(record.RelatedBy) != 0 && len(record.RelatedTo) != 0 && len(record.Related) != 0 {
			answer := []string{}
			nameToAsk := record.RelatedBy
			queryRelated := strings.Join(baseRecord[record.Related], "")

			// if non-related record from specified type/table
			if len(record.Value) != 0 {
				queryRelated = record.Value[0]
			}

			// if record not cached, do it
			if _, ok := relatedRecord[record.Related]; ok == false {
				if record.Multiple {
					relatedRecord[record.Related], err = storage.CurrentStorage.SearchMultiple(
						record.RelatedTo,
						nameToAsk,
						queryRelated,
					)
				} else {
					relatedRecord[record.Related], err = storage.CurrentStorage.SearchRelated(
						record.RelatedTo,
						nameToAsk,
						queryRelated,
					)
				}
			}
			// collects all fields into answer
			for _, item := range record.Name {
				if result, ok := relatedRecord[record.Related][item]; ok {
					answer = append(answer, result...)
				}
			}
			mapp.Fields[index] = mapper.Field{
				Key:      record.Key,
				Value:    answer,
				Format:   record.Format,
				Multiple: record.Multiple,
				Hide:     record.Hide,
			}
		}
	}

	return mapp, nil
}

// prepares join and multiple actions in the answer
func prepareAnswer(mapp *mapper.Record, keys []string) (answer string) {
	for _, index := range keys {
		key := mapp.Fields[index].Key
		if mapp.Fields[index].Hide == true {
			answer = strings.Join([]string{answer, key, "", "\n"}, "")
		} else {
			if mapp.Fields[index].Multiple == true {
				for _, value := range mapp.Fields[index].Value {
					answer = strings.Join([]string{answer, key, value, "\n"}, "")
				}
			} else {
				var value string
				if mapp.Fields[index].Format != "" {
					value = customJoin(mapp.Fields[index].Format, mapp.Fields[index].Value)
				} else {
					value = strings.Join(mapp.Fields[index].Value, " ")
				}
				answer = strings.Join([]string{answer, key, value, "\n"}, "")
			}
		}
	}

	return answer
}

// loads all defined tags from preassigned data before join
func customJoin(format string, value []string) string {
	// template of date to parse
	var templateDateFormat = "2006-01-02 15:04:05"

	for _, item := range value {
		if strings.Contains(format, "{date}") == true {
			buildTime, err := time.Parse(templateDateFormat, item)
			if err != nil && len(strings.TrimSpace(item)) == 0 {
				buildTime = time.Now()
			}
			format = strings.Replace(format, "{date}", buildTime.Format(time.RFC3339), 1)
		}
		format = strings.Replace(format, "{string}", item, 1)
	}
	format = strings.Replace(format, "{string}", "", -1)

	return strings.Trim(format, ". ")
}
