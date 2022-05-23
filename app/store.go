package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

type Store interface {
	Insert(item CheckBoxItem) error
	List() ([]CheckBoxItem, error)
	Close() error
}

const pg = "postgres"
const embed = "embedded"

func createStore(storeSelect, dsn string) (Store, error) {
	if storeSelect == pg {
		return newDBStore(dsn)
	} else if storeSelect == embed {
		return &inMemoryStore{}, nil
	}
	return nil, fmt.Errorf("no valid store selection, valid options are: %q or %q", pg, embed)
}

type inMemoryStore struct {
	m     sync.Mutex
	Items []CheckBoxItem
}

func (s *inMemoryStore) Insert(item CheckBoxItem) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.Items = append(s.Items, item)
	return nil
}

func (s *inMemoryStore) List() ([]CheckBoxItem, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var itemsCopy []CheckBoxItem
	itemsCopy = append(itemsCopy, s.Items...)
	return itemsCopy, nil
}

func (s *inMemoryStore) Close() error {
	return nil
}

type dbStore struct {
	m            sync.Mutex
	db           *sql.DB
	initComplete bool
}

func newDBStore(dsn string) (*dbStore, error) {
	db, err := sql.Open("pgx", dsn)
	return &dbStore{db: db}, err
}

func (s *dbStore) Close() error {
	return s.db.Close()
}

const createTable = `CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, description VARCHAR NOT NULL)`

// Delay database setup until needed for faster application startup
func (s *dbStore) maybeSetup(ctx context.Context) error {
	if s.initComplete { // fast path check
		return nil
	}

	s.m.Lock()
	defer s.m.Unlock()
	// safe check
	if s.initComplete {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	_, err = tx.ExecContext(ctx, createTable)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("table create failed: %v, unable to rollback: %v", err, rollbackErr)
		}
		return fmt.Errorf("table create failed: %v, rollback successful", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing table creation: %v", err)
	}

	s.initComplete = true
	log.Printf("table todos created")
	return nil
}

func (s *dbStore) Insert(item CheckBoxItem) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	if err := s.maybeSetup(ctx); err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, execErr := tx.ExecContext(ctx, "INSERT INTO todos(description) VALUES($1)", item.Description)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("update failed %v, unable to rollback %v", execErr, rollbackErr)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit tx: %v", err)
	}
	return err
}

const selectTodos = "SELECT id,description FROM todos ORDER BY id"

func (s *dbStore) List() ([]CheckBoxItem, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	if err := s.maybeSetup(ctx); err != nil {
		log.Printf("ERROR: %s", err)
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, selectTodos)
	if err != nil {
		return nil, fmt.Errorf("query failed %v", err)
	}
	defer rows.Close()

	var out []CheckBoxItem
	for rows.Next() {
		var item CheckBoxItem
		err := rows.Scan(&item.ID, &item.Description)
		if err != nil {
			return nil, fmt.Errorf("while scanning rows: %v", err)
		}
		out = append(out, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("got err after scanning rows: %v", err)
	}
	return out, nil
}
