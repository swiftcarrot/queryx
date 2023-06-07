package adapter

import (
	"regexp"
)

type Migration struct {
	Path      string
	Version   string
	Name      string
	Direction string
}

var mrx = regexp.MustCompile(`^(\d+)_([^.]+)\.(up|down)\.sql$`)

func ParseMigrationFilename(filename string) (*Migration, error) {
	matches := mrx.FindAllStringSubmatch(filename, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	m := matches[0]

	return &Migration{
		Version:   m[1],
		Name:      m[2],
		Direction: m[3],
	}, nil
}

type Migrations []*Migration
type UpMigrations Migrations
type DownMigrations Migrations

func (ms UpMigrations) Len() int {
	return len(ms)
}

func (ms UpMigrations) Less(i, j int) bool {
	return ms[i].Version < ms[j].Version
}

func (ms UpMigrations) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms DownMigrations) Len() int {
	return len(ms)
}

func (ms DownMigrations) Less(i, j int) bool {
	return ms[i].Version > ms[j].Version
}

func (ms DownMigrations) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
