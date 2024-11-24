INSERT INTO task.task_status (status_name) 
VALUES ('To Do'), ('In Progress'), ('Done') 
ON CONFLICT DO NOTHING;