package category

import "encoding/json"

var dbResult struct {
	ID         int64           `db:"id"`
	Name       string          `db:"name"`
	ParentID   *int64          `db:"parent.id"`
	ParentName *string         `db:"parent.name"`
	Children   json.RawMessage `db:"children"`
}
