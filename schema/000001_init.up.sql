CREATE TABLE groups
(
    number   VARCHAR(255) PRIMARY KEY,
    subjects VARCHAR(255)[]
);

CREATE TABLE users
(
    id            uuid PRIMARY KEY,
    full_name     VARCHAR(255)        NOT NULL,
    group_number  VARCHAR(255),
    username      VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    role          VARCHAR(255)        NOT NULL,
    FOREIGN KEY (group_number) REFERENCES groups (number)
);

CREATE TABLE tasks
(
    id         uuid PRIMARY KEY,
    label      VARCHAR(255) NOT NULL,
    text       VARCHAR(255),
    deadline   TIMESTAMP    NOT NULL,
    points     INTEGER,
    closed     BOOLEAN      NOT NULL,
    teacher_id uuid         NOT NULL,
    subject    VARCHAR(255) NOT NULL,
    is_key_point BOOLEAN NOT NULL,
    file_name  VARCHAR(255),
    student_id uuid         NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP,
    FOREIGN KEY (teacher_id) REFERENCES users (id),
    FOREIGN KEY (student_id) REFERENCES users (id)
);

CREATE TABLE answers
(
    id         uuid PRIMARY KEY,
    message    VARCHAR(255),
    file_name  VARCHAR(255),
    task_id    uuid      NOT NULL,
    points     INTEGER,
    approved   BOOLEAN   not null,
    created_at TIMESTAMP not null,
    FOREIGN KEY (task_id) REFERENCES tasks (id)
);

