# cql2csv

This is a short and sweet program for running queries against a server implementing the Cassandra Query Language and printing out the results in CSV format.

The original reason I wrote this was so I could see a tabular representation of `DESCRIBE` calls. When you run `DESCRIBE SCHEMA`, for instance, `cqlsh` will only show you the values of the `create_statement` column concatenated together. However, if you run that same query with `cql2csv`, you'll get to see the columns being hidden. This was useful when I was testing compatibility between Cassandra and ScyllaDB.

Other than that, this might be useful for shell scripting and things like that. Whatever the case, this is a simple program that I felt like was worth putting in a public repository. If you have ideas how to make this more useful, please fork it and experiment!

## Installation

```sh
go install github.com/decadentsoup/cql2csv
```

## Usage

Command:

```sh
cql2csv 'SELECT * FROM ks.users;'
```

Result:

```csv
handle (varchar),password (blob),user_id (uuid)
robot,[],b71a22bb-696c-4bdb-a6c1-c314ab00c4f3
decadentsoup,[222 173 190 239],2834ec9b-0187-49c6-bb75-39389af49765
```

See `cql2csv --help` for a full list of options.
