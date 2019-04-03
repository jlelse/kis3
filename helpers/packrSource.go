package helpers

// This is a quick hack to get Packr working with golang-migrate

import (
	"bytes"
	"github.com/gobuffalo/packr/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type PackrSource struct {
	lock       sync.Mutex
	Box        *packr.Box
	migrations *source.Migrations
}

func (s *PackrSource) loadMigrations() (*source.Migrations, error) {
	migrations := source.NewMigrations()
	for _, filename := range s.Box.List() {
		migration, err := source.Parse(filename)
		if err != nil {
			return nil, err
		}
		migrations.Append(migration)
	}
	return migrations, nil
}

func (s *PackrSource) Open(url string) (source.Driver, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if migrations, err := s.loadMigrations(); err != nil {
		return nil, err
	} else {
		s.migrations = migrations
		return s, nil
	}
}

func (s *PackrSource) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.migrations = nil
	return nil
}

func (s *PackrSource) First() (version uint, err error) {
	if v, ok := s.migrations.First(); !ok {
		return 0, os.ErrNotExist
	} else {
		return v, nil
	}
}

func (s *PackrSource) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := s.migrations.Prev(version); !ok {
		return 0, os.ErrNotExist
	} else {
		return v, nil
	}
}

func (s *PackrSource) Next(version uint) (nextVersion uint, err error) {
	if v, ok := s.migrations.Next(version); !ok {
		return 0, os.ErrNotExist
	} else {
		return v, nil
	}
}

func (s *PackrSource) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if migration, ok := s.migrations.Up(version); !ok {
		return nil, "", os.ErrNotExist
	} else {
		b, _ := s.Box.Find(migration.Raw)
		return ioutil.NopCloser(bytes.NewBuffer(b)),
			migration.Identifier,
			nil
	}
}

func (s *PackrSource) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if migration, ok := s.migrations.Down(version); !ok {
		return nil, "", migrate.ErrNilVersion
	} else {
		b := s.Box.Bytes(migration.Raw)
		return ioutil.NopCloser(bytes.NewBuffer(b)),
			migration.Identifier,
			nil
	}
}
