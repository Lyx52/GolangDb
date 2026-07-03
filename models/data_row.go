package models

type ColumnValue any

type DataRow struct {
	Values []ColumnValue
}

func NewDataRow(values ...ColumnValue) *DataRow {
	return &DataRow{
		Values: values,
	}
}
