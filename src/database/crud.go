package database

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

// Ping makes a simple ping with 3 second timeout
func Ping() error {
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return instance.DB.PingContext(ctx)
}

// ExecuteQuery executes queries such as INSERT, UPDATE or DELETE
func ExecuteQuery(query string) (sql.Result, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return instance.DB.Unsafe().ExecContext(ctx, query)
}

// ExecuteNamedQuery executes queries such as INSERT, UPDATE or DELETE with named parameters
func ExecuteNamedQuery(query string, arg interface{}) (sql.Result, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return instance.DB.Unsafe().NamedExecContext(ctx, query, arg)
}

// GetSingleRecordNamedQuery selects a single record from a named query and parses it to the destination
func GetSingleRecordNamedQuery(destination interface{}, query string, args interface{}) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	namedStatement, err := instance.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	return namedStatement.Unsafe().GetContext(ctx, destination, args)
}

// GetSingleRecord used for a single record SELECT statements
func GetSingleRecord(destination interface{}, query string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return instance.DB.Unsafe().GetContext(ctx, destination, query)
}

// GetMultipleRecords selects multiple records from the database
func GetMultipleRecords(destination interface{}, query string) (err error) {
	//err = backoff.Retry(func() error {
	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return instance.DB.Unsafe().SelectContext(ctx, destination, query)
	//}, utils2.RetryConfig())
}

// GetMultipleRecordsNamedQuery selects multiple records from the database from a query with named parameters
func GetMultipleRecordsNamedQuery(destination interface{}, query string, input map[string]interface{}) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	parsedQuery, arguments, err := sqlx.Named(query, input)
	if err != nil {
		return nil
	}
	outputQuery, args, err := sqlx.In(parsedQuery, arguments...)
	if err != nil {
		return nil
	}
	outputQuery = instance.DB.Rebind(outputQuery)
	return instance.DB.Unsafe().SelectContext(ctx, destination, outputQuery, args...)
}
