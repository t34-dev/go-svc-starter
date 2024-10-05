CREATE TABLE users
(
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(64) NOT NULL UNIQUE,
    nickname        VARCHAR(64) NOT NULL UNIQUE,
    password_hash   VARCHAR(60) NOT NULL,
    email_confirmed TIMESTAMP NULL,
    logo            VARCHAR(255)         DEFAULT NULL,
    is_block        BOOLEAN     NOT NULL DEFAULT FALSE,
    created         TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN users.email IS 'User''s email';
COMMENT ON COLUMN users.nickname IS 'User''s nickname';
COMMENT ON COLUMN users.password_hash IS 'Password hash';
COMMENT ON COLUMN users.email_confirmed IS 'Deadline for email confirmation';
COMMENT ON COLUMN users.logo IS 'Avatar URL';
COMMENT ON COLUMN users.is_block IS 'Is the user blocked?';
COMMENT ON COLUMN users.created IS 'User creation date';

CREATE TABLE roles
(
    id   SERIAL PRIMARY KEY,
    role VARCHAR(100) NOT NULL UNIQUE
);

COMMENT ON COLUMN roles.id IS 'Role ID';
COMMENT ON COLUMN roles.role IS 'Role name';

-- Assigning roles to users
CREATE TABLE user_roles
(
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles (id) ON DELETE CASCADE
);

COMMENT ON COLUMN user_roles.user_id IS 'Role for user';
COMMENT ON COLUMN user_roles.role_id IS 'Role';

INSERT INTO roles (role)
VALUES ('Admin'),
       ('Vip');

CREATE TABLE user_sessions
(
    user_id      INT REFERENCES users (id) ON DELETE CASCADE,
    refresh      VARCHAR(160) NOT NULL UNIQUE,
    expires_at   TIMESTAMP    NOT NULL,
    finger_print VARCHAR(36)  NOT NULL,
    user_agent   VARCHAR(255) NOT NULL,
    update       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_sessions_encrypt
(
    id          VARCHAR(143) NOT NULL UNIQUE,
    private_key BYTEA        NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN user_sessions.user_id IS 'User for whom the session is created';
COMMENT ON COLUMN user_sessions.refresh IS 'Refresh token';
COMMENT ON COLUMN user_sessions.expires_at IS 'When this token expires';
COMMENT ON COLUMN user_sessions.finger_print IS 'Unique device HASH';
COMMENT ON COLUMN user_sessions.user_agent IS 'Login device';
