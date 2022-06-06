-- +goose Up
ALTER TABLE workspaces ADD COLUMN current_run_id TEXT;
ALTER TABLE workspaces ADD CONSTRAINT workspace_current_run_id_fk FOREIGN KEY (current_run_id) REFERENCES runs ON UPDATE CASCADE;

-- +goose Down
ALTER TABLE workspaces DROP CONSTRAINT workspace_current_run_id_fk;
ALTER TABLE workspaces DROP COLUMN current_run_id;
