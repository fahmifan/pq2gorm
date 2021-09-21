# pq2gorm - Generate [gorm](https://github.com/jinzhu/gorm) model structs from PostgreSQL database schema

[![Build Status](https://travis-ci.org/wantedly/pq2gorm.svg?branch=master)](https://travis-ci.org/wantedly/pq2gorm)

pq2gorm is a generator of [gorm](https://github.com/jinzhu/gorm) model structs from a PostgresSQL database.

* Input: Connection URI of a PostgresSQL database.
* Output: Model definitions based on [gorm](https://github.com/jinzhu/gorm) annotated struct.

## How to build and install

Prepare Go 1.6 or higher.
Go 1.5 is acceptable, but `GO15VENDOREXPERIMENT=1` must be set.

After installing required version of Go, you can build and install `pq2gorm` by

```bash
$ go get -d -u github.com/wantedly/pq2gorm
$ cd $GOPATH/src/github.com/wantedly/pq2gorm
$ make
$ make install
```

`make` generates a binary into `bin/pq2gorm`.
`make install` put it to `$GOPATH/bin`.

## How to use

Run `pq2gorm` with Connection URI of a PostgresSQL database.
Connection URI is necessary for running.

### Usage

```bash
$ pq2gorm
Usage: Generate gorm model structs from PostgreSQL database schema.
  -d string
        Set output path (default "./")
  -dir string
        Set output path (default "./")
  -t string
        Target tables (table1,table2,...) (default: all tables)
  -tables string
        Target tables (table1,table2,...) (default: all tables)
```

**Example 1:** Generate gorm model files of all tables in current directory.

```bash
$ pq2gorm "postgresql://user:password@host:port/dbname?sslmode=disable"
```

For example, user model user.go as shown below will be generated:

```go
type User struct {
    ID uint `json:"id"`
    ...
}
```

**Example 2:** Generate gorm model files of all tables in `./out` directory.

```bash
$ pq2gorm "postgresql://user:password@host:port/dbname?sslmode=disable"  -dbtype postgres -d ./out
```

**Example 3:** Generate gorm model files of `profiles` and `users` tables.

```bash
$ pq2gorm "postgresql://user:password@host:port/dbname?sslmode=disable" -d ./out -t profiles,users
```

**Example 4** Generate gorm model files of all tables in `./out` directory for cockroachdb.

```bash
$ pq2gorm "postgresql://user:password@host:port/dbname?sslmode=disable"  -dbtype cockroachdb -dbname $dbname -d ./out
```

If the directory `./out` does not exist, `pq2gorm` creates `./out` directory with output files.

## License
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
