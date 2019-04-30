-- +migrate Up
ALTER TABLE views
    ADD COLUMN useragent TEXT DEFAULT '' NOT NULL;
