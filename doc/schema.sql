-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2024-12-12T04:24:03.593Z

CREATE TABLE "accounts" (
  "id" bigint PRIMARY KEY NOT NULL,
  "owner" "character varying" NOT NULL,
  "balance" bigint NOT NULL,
  "currency" "character varying" NOT NULL,
  "create_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigint PRIMARY KEY NOT NULL,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "schema_migrations" (
  "version" bigint PRIMARY KEY NOT NULL,
  "dirty" boolean NOT NULL
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY NOT NULL,
  "username" "character varying",
  "refresh_token" "character varying" NOT NULL,
  "user_agent" "character varying" NOT NULL,
  "client_ip" "character varying" NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamp NOT NULL,
  "create_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigint PRIMARY KEY NOT NULL,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "username" "character varying" PRIMARY KEY NOT NULL,
  "hashed_password" "character varying" NOT NULL,
  "full_name" "character varying" NOT NULL,
  "email" "character varying" UNIQUE NOT NULL,
  "password_changed_at" timestamp NOT NULL DEFAULT ('0001-01-01 00:00:00+00'::timestampwithtimezone),
  "create_at" timestamp NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "owner_currency_key" ON "accounts" ("owner", "currency");

CREATE INDEX "accounts_owner_idx" ON "accounts" USING BTREE ("owner");

CREATE INDEX "entries_account_id_idx" ON "entries" USING BTREE ("account_id");

CREATE INDEX "transfers_from_account_id_idx" ON "transfers" USING BTREE ("from_account_id");

CREATE INDEX "transfers_from_account_id_to_account_id_idx" ON "transfers" USING BTREE ("from_account_id", "to_account_id");

CREATE INDEX "transfers_to_account_id_idx" ON "transfers" USING BTREE ("to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "accounts" ADD CONSTRAINT "accounts_owner_fkey" FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD CONSTRAINT "entries_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "sessions" ADD CONSTRAINT "sessions_username_fkey" FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "transfers" ADD CONSTRAINT "transfers_from_account_id_fkey" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD CONSTRAINT "transfers_to_account_id_fkey" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
