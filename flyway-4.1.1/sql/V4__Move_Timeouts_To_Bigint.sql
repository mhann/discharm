ALTER TABLE timeouts
ALTER COLUMN guild_id TYPE bigint;

ALTER TABLE timeouts
ALTER COLUMN target_user_id TYPE bigint;

ALTER TABLE timeouts
ALTER COLUMN creator_user_id TYPE bigint;