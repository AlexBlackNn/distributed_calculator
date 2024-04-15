CREATE TABLE IF NOT EXISTS app_users
(
    uid uuid PRIMARY KEY,
    user_name varchar(40) NOT NULL
);

ALTER TABLE operations ADD user_id uuid;
ALTER TABLE settings ADD user_id uuid UNIQUE;

ALTER TABLE operations
    ADD CONSTRAINT operations_app_users_id_id_fkey
        FOREIGN KEY (user_id) REFERENCES app_users(uid);

ALTER TABLE settings
    ADD CONSTRAINT settings_app_users_id_id_fkey
        FOREIGN KEY (user_id) REFERENCES app_users(uid);
