-- +migrate Up
ALTER TABLE views
    ADD COLUMN os TEXT DEFAULT '' NOT NULL;
