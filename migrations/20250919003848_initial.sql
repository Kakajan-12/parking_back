-- Modify "cameras" table
ALTER TABLE "cameras" ADD COLUMN "channel_name" character varying(255) NULL, ADD COLUMN "channel_token" character varying(255) NULL;
