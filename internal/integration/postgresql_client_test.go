package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swiftcarrot/queryx/internal/integration/db/queryx"
)

func TestVector(t *testing.T) {
	_, err := c.QueryItem().DeleteAll()
	require.NoError(t, err)

	item1, err := c.QueryItem().Create(c.ChangeItem().SetEmbedding([]float32{1, 2, 3}))
	require.NoError(t, err)
	require.Equal(t, item1.Embedding.Val, []float32{1, 2, 3})

	item2, err := c.QueryItem().Create(c.ChangeItem().SetEmbedding([]float32{4, 5, 6}))
	require.NoError(t, err)
	require.Equal(t, item2.Embedding.Val, []float32{4, 5, 6})

	type Foo struct {
		embedding queryx.Vector `db:"embedding"`
	}
	var rows []Foo
	err = c.Query("SELECT embedding FROM items ORDER BY embedding <-> '[3,1,2]'").Scan(&rows)
	require.NoError(t, err)
	require.Equal(t, []Foo{
		{queryx.NewVector([]float32{1, 2, 3})},
		{queryx.NewVector([]float32{4, 5, 6})},
	}, rows)
}
