ALTER TABLE timeouts
ADD COLUMN created timestamp;

ALTER TABLE timeouts
ADD COLUMN length_seconds bigint;

ALTER TABLE timeouts
ADD COLUMN known_expired bool;