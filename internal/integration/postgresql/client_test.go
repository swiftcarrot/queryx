package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swiftcarrot/queryx/internal/integration/postgresql/db"
)

var c *db.QXClient

func init() {
	client, err := db.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	c = client
}

func TestStringArray(t *testing.T) {
	emails := []string{"test1@example.com", "test2@example.com"}
	user, err := c.QueryUser().Create(c.ChangeUser().SetEmails(emails))
	require.NoError(t, err)
	require.Equal(t, emails, user.Emails)
	row, err := c.QueryUser().Find(user.ID)
	require.NoError(t, err)
	require.Equal(t, emails, row.Emails)
}
