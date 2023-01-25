package handlers

import (
	"database/sql"
)

func DBExecute(db sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	cols, _ := rows.Columns()

	var result []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}
	
		row := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			row[colName] = *val
		}
		
		result = append(result, row)
	}

	return result, nil
}