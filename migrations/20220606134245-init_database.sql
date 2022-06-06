-- +migrate Up
CREATE TABLE IF NOT EXISTS "groups" (
     "id" SERIAL PRIMARY KEY,
     "created_at" TIMESTAMPTZ NOT NULL,
     "updated_at" TIMESTAMPTZ NOT NULL,
     "deleted_at" TIMESTAMPTZ,
     "name" VARCHAR(300) NOT NULL,
     "super_group_id" INTEGER REFERENCES "groups" ("id") ON DELETE SET NULL
);
CREATE TABLE IF NOT EXISTS "humans" (
    "id" SERIAL PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ,
    "name" VARCHAR(300) NOT NULL,
    "surname" VARCHAR(300) NOT NULL,
    "birthdate" TIMESTAMPTZ NOT NULL
);
CREATE TABLE IF NOT EXISTS "humans_groups" (
    "human_id" INT NOT NULL REFERENCES "humans" ("id") ON DELETE CASCADE,
    "group_id" INT NOT NULL REFERENCES "groups" ("id") ON DELETE CASCADE
);
-- +migrate Down
DROP TABLE "groups";
DROP TABLE "humans";
DROP TABLE "humans_groups";
