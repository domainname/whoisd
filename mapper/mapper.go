package mapper

// Record - standard record (struct) for mapper package
type Record struct {
	Fields map[string]Field // a list of fields from "01" to last number "nn" in ascending order
}

// Field - representation of one field
type Field struct {
	Key       string   // the label for the field in whois output
	Value     []string // used if the field has prearranged value
	Name      []string // the name of the field in the database, if the field is not prearranged ("value" is not defined)
	Format    string   // special instructions to indicate how to display field
	Multiple  bool     // if this option is set to 'true', then for each value will be repeated label in whois output
	Hide      bool     // if this option is set to 'true', the value of the field will not shown in whois output
	Related   string   // the name of the field in the database through which the request for
	RelatedBy string   // the name of the field in the database through which the related request for
	RelatedTo string   // the name of the table/type in the database through which made a relation
}

// Count - get count of fields
func (mapper *Record) Count() int {
	return len(mapper.Fields)
}
