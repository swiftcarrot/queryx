class User:
    def __init__(self, *, name: str = None, email: str = None):
        self.name = name
        self.email = email

    def __str__(self):
        return f"Name: {self.name}, Email: {self.email}"


user1 = User()
user2 = User(name="john")

print(user1)  # Email: None
print(user2)  # Email: None

dict1 = {"name": "john", "email": None}
print(dict1.get("email") is None)  # True
print(dict1.get("age") is None)  # True


# user = c.queryUser().create(name="john")
