-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table with UUID
CREATE TABLE users (
                       id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email      VARCHAR(64) UNIQUE NOT NULL,
                       username   VARCHAR(64) UNIQUE NOT NULL,
                       password   CHAR(60)           NOT NULL,
                       logo_url   VARCHAR(255)                DEFAULT NULL,
                       is_blocked BOOLEAN            NOT NULL DEFAULT FALSE,
                       created_at TIMESTAMP WITH TIME ZONE    DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE    DEFAULT CURRENT_TIMESTAMP
);

-- Create sessions table with UUID
CREATE TABLE sessions (
                          id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          user_id       UUID REFERENCES users (id) ON DELETE CASCADE,
                          device_key    VARCHAR(255)             NOT NULL,
                          device_name   VARCHAR(255)             NOT NULL,
                          expires_at    TIMESTAMP WITH TIME ZONE NOT NULL,
                          created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          UNIQUE (user_id, device_key)
);

-- Create roles table with regular SERIAL
CREATE TABLE roles (
                       id   SERIAL PRIMARY KEY,
                       role VARCHAR(100) NOT NULL UNIQUE
);

-- Create user_roles table
CREATE TABLE user_roles (
                            user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                            role_id INTEGER NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
                            PRIMARY KEY (user_id, role_id)
);

-- Insert initial roles
INSERT INTO roles (role) VALUES ('Admin'), ('Vip');

-- Create index for faster lookups
CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);

-- Add trigger to update the 'updated_at' column
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

-- Add trigger for users table
CREATE TRIGGER update_user_modtime
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

-- Add trigger for sessions table
CREATE TRIGGER update_session_modtime
    BEFORE UPDATE ON sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();
