CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL UNIQUE,
    description TEXT
);
INSERT INTO roles(name, description, level)
VALUES (
        'user',
        'can create posts and comments',
        1
    );
INSERT INTO roles(name, description, level)
VALUES (
        'moderator',
        'can update other users posts and comments',
        2
    );
INSERT INTO roles(name, description, level)
VALUES (
        'admin',
        'can update and delete posts and comments',
        3
    );