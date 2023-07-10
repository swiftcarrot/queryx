import pytest
import psycopg
import asyncio
import asyncpg
import datetime

url = "postgresql://postgres:postgres@localhost:5432/queryx_test?sslmode=disable"


async def main():
    # Establish a connection to an existing database named "test"
    # as a "postgres" user.
    conn = await asyncpg.connect(url)

    await conn.execute("INSERT INTO tags(name) VALUES($1)", "test", )
    row = await conn.fetchrow('SELECT * FROM tags WHERE name = $1', 'test')
    # *row* now contains
    # asyncpg.Record(id=1, name='Bob', dob=datetime.date(1984, 3, 1))
    print(row['id'])

    # Close the connection.
    await conn.close()

asyncio.get_event_loop().run_until_complete(main())


class User:
    def __init__(self, id, name) -> None:
        self.id = id
        self.name = name

    def __str__(self) -> str:
        return "User(id={}, name={})".format(self.id, self.name)


class SelectStatement:
    def __init__(self) -> None:
        pass


# TODO: convert ? to %s
def rebind(sql, args):
    return sql, args


class QXClient:
    def __init__(self) -> None:
        self.conn = psycopg.connect(url)

    def query_one(self, query, args):
        sql, args = rebind(query, args)
        with self.conn.cursor() as cur:
            cur.execute(query, args)
        return cur.fetchone()

    def query(self, query, args):
        sql, args = rebind(query, args)
        with self.conn.cursor() as cur:
            cur.execute(query, args)
        return cur.fetchall()

    def exec(self, query, args):
        sql, args = rebind(query, args)
        with self.conn.cursor() as cur:
            cur.execute(query, args)

    def query_user(self):
        return UserQuery(self)


class UserChange:
    def __init__(self, input) -> None:
        pass

    def setName(self, name: str):
        self.name = name


class UserQuery:
    def __init__(self, client: QXClient) -> None:
        self.client = client

    def create(self, input):
        return self.client.query_one("insert into users(name) values (%s)", input["name"])


c = QXClient()


def test_query_one():
    user = c.query_user().create({"name": "test"})
    row = c.query_one("select id as user_id from users where id = ?", user.id)
    assert row.user_id == user.id


def test_query():
    user1 = c.query_user().create({"name": "test1"})
    user2 = c.query_user().create({"name": "test2"})
    rows = c.query(
        "select name as user_name from users where id in (?)", [user1.id, user2.id])
    assert rows == [{"user_name": "test1"}, {"user_name", "test2"}]


def test_exec():
    user = c.query_user().create({"name": "test"})
    updated = c.exec("update users set name = ? where id = ?",
                     "test1", user.id)
    assert updated == 1
    deleted = c.exec("delete from users where id = ?", user.id)
    assert deleted == 1


def test_create():
    user = c.query_user().create({"name": "test", type: "admin"})
    assert user.name == "test"
    assert user.type == "admin"
    assert user.id is not None


def test_find():
    user = c.query_user().create({"name": "test"})
    user = c.query_user().find(user.id)
    assert user.name == "test"
    with pytest.raises(Exception):
        c.query_user().find(user.id + 1)


def test_change():
    userChange = UserChange({"name": "test"})
    assert userChange.name == "test"
