package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type DBType string

const (
	DBTypePostgres = DBType("postgres")
	DBTypeCRDB     = DBType("cockroachdb")
)

type DBService interface {
	RetrieveTables(targets []string) ([]string, error)
	RetrieveFields(table string) ([]*Field, error)
	RetrievePrimaryKeys(table string) (map[string]bool, error)
	// Close close connection
	CloseConn() error
}

func main() {
	var (
		dir     string
		ts      string
		dbaname string
		dbtype  string
	)

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %s:
  %s <PostgreSQL URL> [<options>]

Options:
`, os.Args[0], os.Args[0])
		f.PrintDefaults() // Print usage of options
	}
	f.StringVar(&dir, "dir", "./", "Set output path")
	f.StringVar(&dir, "d", "./", "Set output path")
	f.StringVar(&ts, "tables", "", "Target tables (table1,table2,...) (default: all tables)")
	f.StringVar(&ts, "t", "", "Target tables (table1,table2,...) (default: all tables)")
	f.StringVar(&dbaname, "dbname", "", "Database name")
	f.StringVar(&dbtype, "dbtype", "", "Database Type, available type: postgres, cockroachdb")

	f.Parse(os.Args[1:])

	var url string

	for 0 < f.NArg() {
		url = f.Args()[0]
		f.Parse(f.Args()[1:])
	}

	if url == "" {
		f.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Connecting to database...")

	var dbService DBService
	var err error
	switch DBType(dbtype) {
	case DBTypeCRDB:
		dbService, err = NewCRDB(url, dbaname)
	case DBTypePostgres:
		dbService, err = NewPostgres(url)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer dbService.CloseConn()

	var targets []string

	for _, t := range strings.Split(ts, ",") {
		if t != "" {
			targets = append(targets, t)
		}
	}

	tables, err := dbService.RetrieveTables(targets)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	modelParams := map[string]*TemplateParams{}

	for _, table := range tables {
		fmt.Println("Table name: " + table)

		pkeys, err := dbService.RetrievePrimaryKeys(table)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fields, err := dbService.RetrieveFields(table)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		modelParams[table] = GenerateModel(table, pkeys, fields, tables)
	}

	for table, param := range modelParams {
		fmt.Println("Add relation for Table name: " + table)

		AddHasMany(param)

		if err := SaveModel(table, param, dir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
