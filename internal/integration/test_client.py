import pytest
import psycopg


class User:
    def __init__(self, id, name) -> None:
        self.id = id
        self.name = name


class SelectStatement:
    def __init__(self, sql, params) -> None:
        self.sql = sql
        self.params = params


class UserQuery:
    def __init__(self) -> None:
        pass

    def create(self, user):
        return user


class QXClient:
    def __init__(self) -> None:
        url = "postgresql://postgres:postgres@localhost:5432/queryx_test?sslmode=disable"
        self.conn = psycopg.connect(url)

    def query_one(self):
        pass

    def query(self):
        pass

    def exec(self, query, args):
        with self.conn.cursor() as cur:
            cur.execute(query, args)

    def query_user(self):
        return UserQuery()


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
