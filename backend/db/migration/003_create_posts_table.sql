CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userID INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category_id INTEGER NOT NULL,

    FOREIGN KEY (userID) REFERENCES users(id)
    FOREIGN KEY (category_id) REFERENCES categories(id)
);