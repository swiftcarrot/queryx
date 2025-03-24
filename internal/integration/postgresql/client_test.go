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
	strings := []string{"test1", "test2"}
	user, err := c.QueryUser().Create(c.ChangeUser().SetStrings(strings))
	require.NoError(t, err)
	require.Equal(t, strings, user.Strings)
	row, err := c.QueryUser().Find(user.ID)
	require.NoError(t, err)
	require.Equal(t, strings, row.Strings)
}

func TestTextArray(t *testing.T) {
	texts := []string{"test1", "test2"}
	user, err := c.QueryUser().Create(c.ChangeUser().SetTexts(texts))
	require.NoError(t, err)
	require.Equal(t, texts, user.Texts)
	row, err := c.QueryUser().Find(user.ID)
	require.NoError(t, err)
	require.Equal(t, texts, row.Texts)
}

func TestIntegerArray(t *testing.T) {
	integers := []int{1, 2}
	user, err := c.QueryUser().Create(c.ChangeUser().SetIntegers(integers))
	require.NoError(t, err)
	require.Equal(t, integers, user.Integers)
	row, err := c.QueryUser().Find(user.ID)
	require.NoError(t, err)
	require.Equal(t, integers, row.Integers)
}
