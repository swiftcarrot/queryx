package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/swiftcarrot/queryx/internal/integration/db"
	"github.com/swiftcarrot/queryx/internal/integration/db/queryx"
)

var c *db.QXClient

func init() {
	client, err := db.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	c = client
}

func TestQueryOne(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetName("test"))
	require.NoError(t, err)

	var row struct {
		UserID int64 `db:"user_id"`
	}
	err = c.QueryOne("select id as user_id from users where id = ?", user.ID).Scan(&row)
	require.NoError(t, err)
	require.Equal(t, user.ID, row.UserID)
}

func TestQuery(t *testing.T) {
	user1, err := c.QueryUser().Create(c.ChangeUser().SetName("test1"))
	require.NoError(t, err)
	user2, err := c.QueryUser().Create(c.ChangeUser().SetName("test2"))
	require.NoError(t, err)

	type Foo struct {
		UserName string `db:"user_name"`
	}
	var rows []Foo
	err = c.Query("select name as user_name from users where id in (?)", []int64{user1.ID, user2.ID}).Scan(&rows)
	require.NoError(t, err)
	require.Equal(t, []Foo{
		{user1.Name.Val},
		{user2.Name.Val},
	}, rows)
}

func TestExec(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetName("test"))
	require.NoError(t, err)
	updated, err := c.Exec("update users set name = ? where id = ?", "test1", user.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), updated)
	deleted, err := c.Exec("delete from users where id = ?", user.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), deleted)
}

func TestCreate(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetName("user").SetType("admin"))
	require.NoError(t, err)
	require.Equal(t, "user", user.Name.Val)
	require.True(t, user.Name.Valid)
	require.Equal(t, "admin", user.Type.Val)
	require.True(t, user.Type.Valid)
	require.True(t, user.ID > 0)
}

func TestInsertAll(t *testing.T) {
	_, err := c.QueryUserPost().DeleteAll()
	require.NoError(t, err)

	_, err = c.QueryPost().InsertAll(nil)
	require.ErrorIs(t, err, db.ErrInsertAllEmptyChanges)
	_, err = c.QueryPost().InsertAll([]*queryx.PostChange{})
	require.ErrorIs(t, err, db.ErrInsertAllEmptyChanges)

	inserted, err := c.QueryPost().InsertAll([]*queryx.PostChange{
		c.ChangePost().SetTitle("title1"),
		c.ChangePost().SetTitle("title2"),
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), inserted)
}

func TestCreateEmpty(t *testing.T) {
	tag, err := c.QueryTag().Create(nil)
	require.NoError(t, err)
	require.True(t, tag.ID > 0)
	require.False(t, tag.Name.Valid)

	tag, err = c.QueryTag().Create(c.ChangeTag())
	require.NoError(t, err)
	require.True(t, tag.ID > 0)
	require.False(t, tag.Name.Valid)
}

func TestFind(t *testing.T) {
	_, err := c.QueryTag().DeleteAll()
	require.NoError(t, err)

	tag, err := c.QueryTag().Create(c.ChangeTag().SetName("test"))
	require.NoError(t, err)
	tag, err = c.QueryTag().Find(tag.ID)
	require.NoError(t, err)
	require.Equal(t, "test", tag.Name.Val)

	tag, err = c.QueryTag().Find(tag.ID + 1)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Nil(t, tag)
}

func TestOrder(t *testing.T) {
	_, err := c.QueryTag().DeleteAll()
	require.NoError(t, err)
	tag1, err := c.QueryTag().Create(c.ChangeTag().SetName("tag1"))
	require.NoError(t, err)
	tag2, err := c.QueryTag().Create(c.ChangeTag().SetName("tag2"))
	require.NoError(t, err)
	tags, err := c.QueryTag().Order("id desc").All()
	require.NoError(t, err)
	require.Equal(t, []*db.Tag{tag2, tag1}, tags)
}

func TestFirst(t *testing.T) {
	_, err := c.QueryTag().DeleteAll()
	require.NoError(t, err)

	tag1, err := c.QueryTag().Create(c.ChangeTag().SetName("test"))
	require.NoError(t, err)
	tag2, err := c.QueryTag().First()
	require.NoError(t, err)
	require.Equal(t, "test", tag2.Name.Val)
	require.Equal(t, tag1.ID, tag2.ID)

	_, err = c.QueryTag().DeleteAll()
	require.NoError(t, err)

	tag3, err := c.QueryTag().First()
	require.NoError(t, err)
	require.Nil(t, tag3)
}

func TestTime(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetTime("12:10:09"))
	require.NoError(t, err)
	require.Equal(t, "12:10:09", user.Time.Val.Format("15:04:05"))
}

func TestDate(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetDate("2012-11-10"))
	require.NoError(t, err)
	require.Equal(t, "2012-11-10", user.Date.Val.Format("2006-01-02"))

	user, err = c.QueryUser().Where(c.UserDate.EQ("2012-11-10")).First()
	require.NoError(t, err)
	require.Equal(t, "2012-11-10", user.Date.Val.Format("2006-01-02"))
}

func TestDatetime(t *testing.T) {
	s1 := "2012-11-10 09:08:07"
	user, err := c.QueryUser().Create(c.ChangeUser().SetDatetime(s1))
	require.NoError(t, err)
	require.Equal(t, s1, user.Datetime.Val.Format("2006-01-02 15:04:05"))

	user, err = c.QueryUser().Where(c.UserID.EQ(user.ID).And(c.UserDatetime.GE(s1)).And(c.UserDatetime.LE(s1))).First()
	require.NoError(t, err)
	require.Equal(t, s1, user.Datetime.Val.Format("2006-01-02 15:04:05"))

	s2 := "2012-11-10 09:08:07.654"
	user, err = c.QueryUser().Create(c.ChangeUser().SetDatetime(s2))
	require.NoError(t, err)
	require.Equal(t, s2, user.Datetime.Val.Format("2006-01-02 15:04:05.000"))
}

func TestTimestamps(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser())
	require.NoError(t, err)
	require.True(t, user.CreatedAt.Valid)
	require.True(t, user.UpdatedAt.Valid)
	require.True(t, user.CreatedAt.Val.Equal(user.UpdatedAt.Val))

	time.Sleep(time.Millisecond)

	err = user.Update(c.ChangeUser().SetName("new name"))
	require.NoError(t, err)
	require.True(t, user.UpdatedAt.Val.After(user.CreatedAt.Val))
}

func TestUUID(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser())
	require.NoError(t, err)
	require.False(t, user.UUID.Valid)

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
	require.False(t, user.Name.Valid)

	user, _ = c.QueryUser().Find(user.ID)
	require.False(t, user.Name.Valid)
}

func TestJSON(t *testing.T) {
	payload := make(map[string]interface{})
	payload["theme"] = "dark"
	payload["height"] = 170
	payload["weight"] = 65
	user, err := c.QueryUser().Create(c.ChangeUser().SetPayload(payload))
	require.NoError(t, err)
	user, err = c.QueryUser().Find(user.ID)
	require.NoError(t, err)
	require.Equal(t, payload["theme"], user.Payload.Val["theme"])
	// numbers are unmarshalled into float64 by default
	require.Equal(t, float64(payload["height"].(int)), user.Payload.Val["height"])
	require.Equal(t, float64(payload["weight"].(int)), user.Payload.Val["weight"])
}

func TestCompositePrimaryKey(t *testing.T) {
	_, _ = c.QueryCode().DeleteAll()
	code, err := c.QueryCode().Create(c.ChangeCode().SetType("type").SetKey("key"))
	require.NoError(t, err)
	require.Equal(t, "type", code.Type)
	require.Equal(t, "key", code.Key)

	_, err = c.QueryCode().Create(c.ChangeCode().SetType("type").SetKey("key"))
	require.Error(t, err)

	code, err = c.QueryCode().Find("type", "key")
	require.NoError(t, err)
	require.Equal(t, "type", code.Type)
	require.Equal(t, "key", code.Key)
}

func TestWithoutPrimaryKey(t *testing.T) {
	// TODO: test missing Delete() and Find() on ClientQuery and Client
	client, err := c.QueryClient().Create(c.ChangeClient().SetName("client"))
	require.NoError(t, err)
	require.Equal(t, "client", client.Name.Val)
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

	_, err = c.QueryClient().Create(c.ChangeClient().SetName("client"))
	require.NoError(t, err)
	exists, err = c.QueryClient().Exists()
	require.NoError(t, err)
	require.True(t, exists)
}

func TestBelongsTo(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser())
	require.NoError(t, err)
	post, err := c.QueryPost().Create(c.ChangePost().SetAuthorID(user.ID))
	require.NoError(t, err)
	account, err := c.QueryAccount().Create(c.ChangeAccount().SetUserID(user.ID))
	require.NoError(t, err)

	post, err = c.QueryPost().PreloadAuthor().Find(post.ID)
	require.NoError(t, err)
	require.Equal(t, user, post.Author)

	account, err = c.QueryAccount().PreloadUser().Find(account.ID)
	require.NoError(t, err)
	require.Equal(t, user, account.User)
}

func TestAllEmpty(t *testing.T) {
	_, err := c.QueryUser().DeleteAll()
	require.NoError(t, err)

	users, err := c.QueryUser().All()
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, 0, len(users))
}

func TestInEmpty(t *testing.T) {
	_, err := c.QueryUser().DeleteAll()
	require.NoError(t, err)
	users, err := c.QueryUser().Where(c.UserID.In([]int64{})).All()
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, 0, len(users))

	users, err = c.QueryUser().Where(c.UserID.In([]int64{}).And(c.UserID.EQ(1)).And(c.UserID.In([]int64{1}))).All()
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, 0, len(users))
}

func TestWhere(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser().SetName("name").SetType("type"))
	require.NoError(t, err)
	users, err := c.QueryUser().Where(c.UserID.EQ(user.ID)).Where(c.UserName.EQ("name"), c.UserType.EQ("type")).All()
	require.NoError(t, err)
	require.Equal(t, 1, len(users))
	require.Equal(t, user, users[0])

	users, err = c.QueryUser().Where(c.Raw("name = ? and type = ?", "name", "type")).All()
	require.NoError(t, err)
	require.Equal(t, 1, len(users))
	require.Equal(t, user, users[0])
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

func TestHasMany(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser())
	require.NoError(t, err)
	post1, err := c.QueryPost().Create(c.ChangePost())
	require.NoError(t, err)
	userPost1, err := c.QueryUserPost().Create(c.ChangeUserPost().SetUserID(user.ID).SetPostID(post1.ID))
	require.NoError(t, err)

	_, err = c.QueryPost().Create(c.ChangePost())
	require.NoError(t, err)
	_, err = c.QueryUserPost().Create(c.ChangeUserPost())
	require.NoError(t, err)

	userPosts, err := user.QueryUserPosts().All()
	require.NoError(t, err)
	require.Equal(t, []*db.UserPost{userPost1}, userPosts)

	posts, err := user.QueryPosts().All()
	require.NoError(t, err)
	require.Equal(t, []*db.Post{post1}, posts)
}

func TestHasOne(t *testing.T) {
	user, err := c.QueryUser().Create(c.ChangeUser())
	require.NoError(t, err)
	account, err := c.QueryAccount().Create(c.ChangeAccount().SetUserID(user.ID))
	require.NoError(t, err)

	user, err = c.QueryUser().PreloadAccount().Find(user.ID)
	require.NoError(t, err)
	require.Equal(t, account, user.Account)
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

	user, err := c.QueryUser().Find(user1.ID)
	require.NoError(t, err)
	err = user.PreloadUserPosts()
	require.NoError(t, err)
	require.Equal(t, 2, len(user.UserPosts))
	err = user.PreloadPosts()
	require.NoError(t, err)
	require.Equal(t, 2, len(user.Posts))
	err = user.PreloadAccount()
	require.NoError(t, err)
	require.Equal(t, account1.ID, user.Account.ID)
}

func TestTransaction(t *testing.T) {
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

func TestTransactionBlock(t *testing.T) {
	_, err := c.QueryTag().DeleteAll()
	require.NoError(t, err)

	err = c.Transaction(func(tx *db.Tx) error {
		tag1, err := tx.QueryTag().Create(tx.ChangeTag().SetName("tag1"))
		require.NoError(t, err)
		require.Equal(t, "tag1", tag1.Name.Val)
		tag2, err := tx.QueryTag().Create(tx.ChangeTag().SetName("tag2"))
		require.Equal(t, "tag2", tag2.Name.Val)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	total, err := c.QueryTag().Count()
	require.NoError(t, err)
	require.Equal(t, int64(2), total)

	err = c.Transaction(func(tx *db.Tx) error {
		_, err = tx.QueryTag().Create(c.ChangeTag().SetName("tag3"))
		_, err = tx.QueryTag().Create(c.ChangeTag().SetName("tag3"))
		return err
	})
	require.Error(t, err)
	total, err = c.QueryTag().Count()
	require.NoError(t, err)
	require.Equal(t, int64(2), total)
}

func TestChangeJSON(t *testing.T) {
	s := `{"name":"user name","isAdmin":false}`

	userChange := c.ChangeUser()
	err := json.Unmarshal([]byte(s), userChange)
	require.NoError(t, err)
	require.Equal(t, "user name", userChange.Name.Val)
	require.False(t, userChange.IsAdmin.Val)
	require.True(t, userChange.Name.Set)
	require.True(t, userChange.IsAdmin.Set)
	require.False(t, userChange.Age.Set)
}

func TestModelJSON(t *testing.T) {
	tag, err := c.QueryTag().Create(c.ChangeTag().SetName("test"))
	require.NoError(t, err)
	b, err := json.Marshal(tag)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(`{"id":%d,"name":"test"}`, tag.ID), string(b))
}

func TestModelStringer(t *testing.T) {
	client, err := c.QueryClient().Create(c.ChangeClient().SetName("test"))
	require.NoError(t, err)
	require.Equal(t, `(Client name: "test", float: null)`, client.String())

	client, err = c.QueryClient().Create(c.ChangeClient().SetName("test").SetFloat(51.1234))
	require.NoError(t, err)
	require.Equal(t, `(Client name: "test", float: 51.1234)`, client.String())
}
