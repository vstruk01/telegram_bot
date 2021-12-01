CREATE TABLE IF NOT EXISTS Users (
    id          uuid PRIMARY KEY,
    telegram_id INTEGER NOT NULL
);

ALTER TABLE ONLY Users ALTER COLUMN first_name SET DEFAULT '';
ALTER TABLE ONLY Users ALTER COLUMN name SET DEFAULT '';

CREATE TABLE IF NOT EXISTS  Chats (
    id uuid PRIMARY KEY,
    telegram_id INTEGER NOT NULL,
    type TEXT,
    first_name TEXT NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Words (
    id uuid PRIMARY KEY,
    creator uuid,
    word TEXT NOT NULL,
    translate TEXT NOT NULL,
    rate INTEGER NOT NULL,
    translated_count INTEGER NOT NULL
);

ALTER TABLE Words ALTER COLUMN rate SET DEFAULT 0;
ALTER TABLE Words ALTER COLUMN translated_count SET DEFAULT 0;
ALTER TABLE Words ADD CONSTRAINT constraint_creator_id FOREIGN KEY (creator) REFERENCES Users(id);
ALTER TABLE UsersWords ADD CONSTRAINT constraint_user_id FOREIGN KEY (user_id) REFERENCES Users(id);
ALTER TABLE UsersWords ADD CONSTRAINT constraint_word_id FOREIGN KEY (word_id) REFERENCES Words(id);


