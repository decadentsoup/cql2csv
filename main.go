package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"

	"github.com/gocql/gocql"
)

func main() {
	cluster := gocql.NewCluster("127.0.0.1:19042")
	query := "DESCRIBE KEYSPACES"

	writer := csv.NewWriter(os.Stdout)

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	iter := session.Query(query).Iter()
	defer func() {
		if err := iter.Close(); err != nil {
			panic(err)
		}
	}()

	writeHeader(writer, iter.Columns())
	writeRows(writer, iter)

	writer.Flush()

	if err := writer.Error(); err != nil {
		panic(err)
	}
}

func writeHeader(writer *csv.Writer, columns []gocql.ColumnInfo) {
	row := make([]string, len(columns))

	for i, column := range columns {
		row[i] = column.String()
	}

	if err := writer.Write(row); err != nil {
		panic(err)
	}
}

func writeRows(writer *csv.Writer, iter *gocql.Iter) {
	for {
		rowData, err := iter.RowData()
		if err != nil {
			panic(err)
		}

		if !iter.Scan(rowData.Values...) {
			break
		}

		row := make([]string, len(rowData.Values))

		for i, value := range rowData.Values {
			row[i] = fmt.Sprint(reflect.Indirect(reflect.ValueOf(value)).Interface())
		}

		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}
}
