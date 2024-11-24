USE todo_app;
CREATE SCHEMA IF NOT EXISTS task;

CREATE TABLE IF NOT EXISTS task.task_status (
    task_status_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    status_name VARCHAR(255) UNIQUE NOT NULL
);

INSERT IGNORE INTO task.task_status (status_name) 
VALUES ('To Do'), ('In Progress'), ('Done');

CREATE TABLE IF NOT EXISTS task.task (
    task_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_status_id BIGINT NOT NULL,
    FOREIGN KEY (task_status_id) REFERENCES task.task_status(task_status_id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task.task_status_history (
    task_status_history_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL,
    FOREIGN KEY (task_id) REFERENCES task.task(task_id),
    task_status_id BIGINT NOT NULL,
    FOREIGN KEY (task_status_id) REFERENCES task.task_status(task_status_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);