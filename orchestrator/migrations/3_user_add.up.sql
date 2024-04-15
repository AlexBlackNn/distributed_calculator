CREATE TABLE IF NOT EXISTS app_user
(
    uid uuid PRIMARY KEY,
    user_name varchar(40) NOT NULL
);

ALTER TABLE operations ADD user_id uuid;
ALTER TABLE settings ADD user_id uuid;
