-- Create "cameras" table
CREATE TABLE "cameras" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "type" character varying(50) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_cameras_deleted_at" to table: "cameras"
CREATE INDEX "idx_cameras_deleted_at" ON "cameras" ("deleted_at");
-- Create "mac_users" table
CREATE TABLE "mac_users" (
  "id" bigserial NOT NULL,
  "mac_username" text NOT NULL,
  "mac_password" text NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- Create "tariffs" table
CREATE TABLE "tariffs" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "duration" bigint NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "price_amount" text NULL,
  "currency" character varying(20) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_tariffs_duration" to table: "tariffs"
CREATE UNIQUE INDEX "idx_tariffs_duration" ON "tariffs" ("duration");
-- Create "cars" table
CREATE TABLE "cars" (
  "id" bigserial NOT NULL,
  "car_number" character varying(255) NOT NULL,
  "owner_name" character varying(255) NULL,
  "is_staff" boolean NULL DEFAULT false,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_cars_car_number" to table: "cars"
CREATE UNIQUE INDEX "idx_cars_car_number" ON "cars" ("car_number");
-- Create "car_sessions" table
CREATE TABLE "car_sessions" (
  "id" bigserial NOT NULL,
  "start_time" timestamptz NULL,
  "end_time" timestamptz NULL,
  "total_payment_amount" numeric(20,8) NULL DEFAULT 0.0,
  "currency" character varying(20) NOT NULL,
  "status" character varying(100) NULL,
  "reason" text NULL,
  "image_url" text NULL,
  "car_park" character varying(100) NULL,
  "duration" numeric(20,8) NULL DEFAULT 0.0,
  "is_paid" boolean NOT NULL DEFAULT false,
  "camera_token" text NULL,
  "car_id" bigint NOT NULL,
  "camera_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_car_sessions_camera" FOREIGN KEY ("camera_id") REFERENCES "cameras" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_car_sessions_car" FOREIGN KEY ("car_id") REFERENCES "cars" ("id") ON UPDATE CASCADE ON DELETE RESTRICT
);
-- Create "car_subscriptions" table
CREATE TABLE "car_subscriptions" (
  "id" bigserial NOT NULL,
  "total_payment_amount" numeric(20,8) NULL DEFAULT 0.0,
  "currency" character varying(20) NOT NULL,
  "is_paid" boolean NOT NULL DEFAULT false,
  "start_time" timestamptz NULL,
  "end_time" timestamptz NULL,
  "is_active" boolean NULL DEFAULT true,
  "revoked_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "car_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_car_subscriptions_car" FOREIGN KEY ("car_id") REFERENCES "cameras" ("id") ON UPDATE CASCADE ON DELETE RESTRICT
);
-- Create index "idx_car_subscriptions_revoked_at" to table: "car_subscriptions"
CREATE INDEX "idx_car_subscriptions_revoked_at" ON "car_subscriptions" ("revoked_at");
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "username" character varying(100) NOT NULL,
  "full_name" character varying(100) NOT NULL,
  "password" character varying(255) NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "is_superuser" boolean NULL DEFAULT false,
  "role" character varying(20) NOT NULL,
  "car_park" character varying(20) NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
-- Create "user_sessions" table
CREATE TABLE "user_sessions" (
  "id" uuid NOT NULL,
  "ip_address" inet NULL,
  "user_agent" text NULL,
  "expire_at" timestamptz NOT NULL,
  "revoked_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_sessions_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_user_sessions_expire_at" to table: "user_sessions"
CREATE INDEX "idx_user_sessions_expire_at" ON "user_sessions" ("expire_at");
-- Create index "idx_user_sessions_revoked_at" to table: "user_sessions"
CREATE INDEX "idx_user_sessions_revoked_at" ON "user_sessions" ("revoked_at");
-- Create index "idx_user_sessions_user_id" to table: "user_sessions"
CREATE INDEX "idx_user_sessions_user_id" ON "user_sessions" ("user_id");
