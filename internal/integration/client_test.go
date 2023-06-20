package main

import (
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swiftcarrot/queryx/internal/integration/db"
	"github.com/swiftcarrot/queryx/internal/integration/db/queryx"
)

var c *db.QXClient

// TODO: add api to readme
func TestQueryOne(t *testing.T) {
	_, err := c.QueryUser().DeleteAll()
	require.NoError(t, err)
	var changes []*queryx.UserChange
	for i := 0; i < 3; i++ {
		changes = append(changes, c.ChangeUser().SetName("name"+strconv.Itoa(i)))
	}
	res, err := c.QueryUser().BulkCreate(changes)
	require.NoError(t, err)
	require.True(t, res > 0)
	var count struct {
		Count int64 `db:"count"`
	}
	err = c.QueryOne("select count('*') from users").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, int64(3), count.Count)
}
func TestQuery(t *testing.T) {
	_, err := c.QueryUser().DeleteAll()
	require.NoError(t, err)
	var changes []*queryx.UserChange
	for i := 0; i < 3; i++ {
		changes = append(changes, c.ChangeUser().SetName("name"+strconv.Itoa(i)))
	}
	res, err := c.QueryUser().BulkCreate(changes)
	require.NoError(t, err)
	require.True(t, res > 0)

	var users []db.User
	err = c.Query("select * from users where id>?", 0).Scan(&users)
	require.NoError(t, err)
	require.Equal(t, 3, len(users))
}

func TestExec(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetName("has_one"))
	require.NoError(t, err)
	res, err := c.Exec("update users set name='test' where id= ?", user.ID)
	require.NoError(t, err)
	require.True(t, res > 0)
	res, err = c.Exec("delete from users where id= ?", user.ID)
	require.NoError(t, err)
	require.True(t, res > 0)
}

func TestCreate(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetName("user").SetType("admin"))
	require.NoError(t, err)
	require.Equal(t, "user", user.Name.Val)
	require.False(t, user.Name.Null)
	require.Equal(t, "admin", user.Type.Val)
	require.False(t, user.Type.Null)
	require.True(t, user.ID > 0)
}

func TestTime(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetTime("12:12:12"))
	require.NoError(t, err)
	require.Equal(t, "12:12:12", user.Time.Val.Format("15:04:05"))
}

func TestDate(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetDate("2012-12-12"))
	require.NoError(t, err)
	require.Equal(t, "2012-12-12", user.Date.Val.Format("2006-01-02"))
}

func TestDatetime(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetDatetime("2012-12-12 12:12:12"))
	require.NoError(t, err)
	require.Equal(t, "2012-12-12 12:12:12", user.Datetime.Val.Format("2006-01-02 15:04:05"))
}

func TestUUID(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser())
	require.NoError(t, err)
	require.True(t, user.UUID.Null)

	uuid1 := "c7e5b9af-0499-4eca-a7e6-77e10d56987b"
	err = user.Update(c.ChangeUser().SetUUID(uuid1))
	require.NoError(t, err)
	require.Equal(t, uuid1, user.UUID.Val)

	_, err = c.QueryDevice().DeleteAll()
	require.NoError(t, err)

	uuid2 := "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef"
	device, err := c.QueryDevice().Create(c.ChangeDevice().SetID(uuid2))
	require.NoError(t, err)
	require.Equal(t, uuid2, device.ID)

	device, err = c.QueryDevice().Find(uuid2)
	require.NoError(t, err)
	require.Equal(t, uuid2, device.ID)
}

func TestSetNullable(t *testing.T) {
	name := "user"
	user, _ := c.QueryUser().Create(c.ChangeUser().SetNullableName(&name))
	require.Equal(t, name, user.Name.Val)

	user.Update(c.ChangeUser().SetNullableName(nil))
	require.True(t, user.Name.Null)

	user, _ = c.QueryUser().Find(user.ID)
	require.True(t, user.Name.Null)
}

func TestJSON(t *testing.T) {
	payload := make(map[string]interface{})
	payload["theme"] = "dark"
	payload["height"] = 170
	payload["weight"] = 65
	user, _ := c.QueryUser().Create(c.ChangeUser().SetPayload(payload))
	require.Equal(t, payload["theme"], user.Payload.Val["theme"])
	// numbers are unmarshalled into float64 by default
	require.Equal(t, float64(payload["height"].(int)), user.Payload.Val["height"])
	require.Equal(t, float64(payload["weight"].(int)), user.Payload.Val["weight"])
}

func TestPrimaryKey(t *testing.T) {
	_, _ = c.QueryCode().DeleteAll()
	code, err := c.QueryCode().Create(c.ChangeCode().SetType("type").SetKey("key"))
	require.Equal(t, "type", code.Type)
	require.Equal(t, "key", code.Key)
	require.NoError(t, err)

	_, err = c.QueryCode().Create(c.ChangeCode().SetType("type").SetKey("key"))
	require.Error(t, err)

	code, err = c.QueryCode().Find("type", "key")
	require.NoError(t, err)
	require.Equal(t, "type", code.Type)
	require.Equal(t, "key", code.Key)

	_, _ = c.QueryClient().DeleteAll()
	client, err := c.QueryClient().Create(c.ChangeClient().SetName("client"))
	require.NoError(t, err)
	require.Equal(t, "client", client.Name)

	i, _ := c.QueryClient().Delete("client")
	require.Equal(t, int64(1), i)
}

func TestBoolean(t *testing.T) {
	user, _ := c.QueryUser().Create(c.ChangeUser().SetIsAdmin(true))
	require.Equal(t, true, user.IsAdmin.Val)
}

func TestExists(t *testing.T) {
	_, _ = c.QueryClient().DeleteAll()
	exists, err := c.QueryClient().Exists()
	require.NoError(t, err)
	require.False(t, exists)

	_, _ = c.QueryClient().Create(c.ChangeClient().SetName("client"))
	exists, err = c.QueryClient().Exists()
	require.NoError(t, err)
	require.True(t, exists)
}

func TestBelongsTo(t *testing.T) {
	author, err := c.QueryUser().Create(c.ChangeUser().SetName("author"))
	require.NoError(t, err)
	post, err := c.QueryPost().Create(c.ChangePost().SetTitle("post title").SetAuthorID(author.ID))
	require.NoError(t, err)
	post, err = c.QueryPost().PreloadAuthor().Find(post.ID)
	require.NoError(t, err)
	require.Equal(t, author.ID, post.Author.ID)
}

func TestAllEmpty(t *testing.T) {
	_, err := c.QueryUser().DeleteAll()
	require.NoError(t, err)

	users, err := c.QueryUser().All()
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, 0, len(users))
}

func TestInEmptySlice(t *testing.T) {
	_, err := c.QueryUser().DeleteAll()
	require.NoError(t, err)
	users, err := c.QueryUser().Where(c.UserID.In([]int64{})).All()
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, 0, len(users))
}

func TestHasManyEmpty(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetName("user"))
	require.NoError(t, err)
	require.Nil(t, user.UserPosts)
	require.Nil(t, user.Posts)

	user, err = c.QueryUser().PreloadUserPosts().Find(user.ID)
	require.NoError(t, err)
	require.NotNil(t, user.UserPosts)
	require.Equal(t, 0, len(user.UserPosts))

	user, err = c.QueryUser().PreloadPosts().Find(user.ID)
	require.NoError(t, err)
	require.NotNil(t, user.Posts)
	require.NotNil(t, user.UserPosts)
	require.Equal(t, 0, len(user.Posts))
	require.Equal(t, 0, len(user.UserPosts))
}

func TestHasOne(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetName("has_one"))
	require.NoError(t, err)
	account, err := c.QueryAccount().Create(c.ChangeAccount().SetName("account").SetUserID(user.ID))
	require.NoError(t, err)

	user, err = c.QueryUser().PreloadAccount().Find(user.ID)
	require.NoError(t, err)
	require.Equal(t, account.Name, user.Account.Name)
}

func TestPreload(t *testing.T) {
	user1, _ := c.QueryUser().Create(c.ChangeUser().SetName("user1"))
	post1, _ := c.QueryPost().Create(c.ChangePost().SetTitle("post1"))
	post2, _ := c.QueryPost().Create(c.ChangePost().SetTitle("post2"))
	account1, _ := c.QueryAccount().Create(c.ChangeAccount().SetName("account1").SetUserID(user1.ID))

	userPost1, _ := c.QueryUserPost().Create(c.ChangeUserPost().SetUserID(user1.ID).SetPostID(post1.ID))
	userPost2, _ := c.QueryUserPost().Create(c.ChangeUserPost().SetUserID(user1.ID).SetPostID(post2.ID))

	user, _ := c.QueryUser().PreloadPosts().PreloadAccount().Find(user1.ID)
	require.Equal(t, account1.ID, user.Account.ID)

	require.Equal(t, 2, len(user.UserPosts))
	require.Equal(t, userPost1.ID, user.UserPosts[0].ID)
	require.Equal(t, userPost2.ID, user.UserPosts[1].ID)

	require.Equal(t, 2, len(user.Posts))
	require.Equal(t, post1.ID, user.Posts[0].ID)
	require.Equal(t, post2.ID, user.Posts[1].ID)

	post, _ := c.QueryPost().PreloadUserPosts().Find(post1.ID)
	require.Equal(t, 1, len(post.UserPosts))
	require.Equal(t, userPost1.ID, post.UserPosts[0].ID)
}

func TestTx(t *testing.T) {
	tag1, _ := c.QueryTag().Create(c.ChangeTag().SetName("tag1"))
	require.Equal(t, "tag1", tag1.Name.Val)

	total1, _ := c.QueryTag().Count()
	tx, _ := c.Tx()

	tag1, _ = tx.QueryTag().Find(tag1.ID)
	tag1.Update(tx.ChangeTag().SetName("tag1-updated"))

	tx.QueryTag().Create(tx.ChangeTag().SetName("tag2"))
	tx.QueryTag().Create(tx.ChangeTag().SetName("tag3"))

	total2, _ := c.QueryTag().Count()
	require.Equal(t, total1, total2)

	total3, _ := tx.QueryTag().Count()
	require.Equal(t, total1+2, total3)

	tag1, _ = c.QueryTag().Find(tag1.ID)
	require.Equal(t, "tag1", tag1.Name.Val)

	tx.Commit()

	total4, _ := c.QueryTag().Count()
	require.Equal(t, total1+2, total4)

	tag1, _ = c.QueryTag().Find(tag1.ID)
	require.Equal(t, "tag1-updated", tag1.Name.Val)
}

func init() {
	client, err := db.NewClientWithEnv("test")
	if err != nil {
		log.Fatal(err)
	}
	c = client
}
