-- +migrate Up
ALTER TABLE views
    ADD COLUMN useragent TEXT DEFAULT '' NOT NULL;

-- +migrate Down
BEGIN TRANSACTION;
CREATE TABLE views_new AS
SELECT url, time, ref
FROM views;
DROP TABLE views;
ALTER TABLE views_new
    RENAME TO views;
COMMIT;
