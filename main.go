package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/gocql/gocql"
)

type cli struct {
	Host        []string `kong:"short='H',default='localhost:9042',help='Add a host address to try when connecting to the cluster.'"` //nolint:lll
	Keyspace    string   `kong:"short='k',help='Set default keyspace.'"`
	Consistency string   `kong:"short='c',default='ONE',help='Set consistency level to use for commands.'"`
	Query       string   `kong:"arg,help='Query to execute against the cluster.'"`
}

func main() {
	var cli cli

	kong.Parse(&cli)

	writer := csv.NewWriter(os.Stdout)

	cluster := gocql.NewCluster(cli.Host...)
	cluster.Keyspace = cli.Keyspace
	cluster.Consistency = gocql.ParseConsistency(cli.Consistency)

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	iter := session.Query(cli.Query).Iter()
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
		row[i] = fmt.Sprintf("%v (%v)", column.Name, column.TypeInfo)
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
