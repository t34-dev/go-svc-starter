CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(64) UNIQUE NOT NULL,
    username   VARCHAR(64) UNIQUE NOT NULL,
    password   VARCHAR(128)       NOT NULL,
    logo_url   VARCHAR(255)                DEFAULT NULL,
    is_blocked BOOLEAN            NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE    DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE    DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (email, username)
);

CREATE TABLE devices
(
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER REFERENCES users (id),
    device_key    VARCHAR(255)             NOT NULL,
    device_name   VARCHAR(255)             NOT NULL,
    last_used     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    refresh_token VARCHAR(255)             NOT NULL,
    expires_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, device_key)
);

CREATE TABLE roles
(
    id   SERIAL PRIMARY KEY,
    role VARCHAR(100) NOT NULL UNIQUE
);
CREATE TABLE user_roles
(
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles (id) ON DELETE CASCADE
);

INSERT INTO roles (role)
VALUES ('Admin'),
       ('Vip');
