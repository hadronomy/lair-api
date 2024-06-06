-- Create "lairs" table
CREATE TABLE "public"."lairs" (
  "id" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "name" text NULL,
  "owner" text NULL,
  "private" boolean NULL,
  PRIMARY KEY ("id")
);
-- Create "models" table
CREATE TABLE "public"."models" (
  "id" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "minions" table
CREATE TABLE "public"."minions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "name" text NULL,
  "level" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create "lair_minions" table
CREATE TABLE "public"."lair_minions" (
  "lair_id" text NOT NULL,
  "minion_id" bigint NOT NULL,
  PRIMARY KEY ("lair_id", "minion_id"),
  CONSTRAINT "fk_lair_minions_lair" FOREIGN KEY ("lair_id") REFERENCES "public"."lairs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_lair_minions_minion" FOREIGN KEY ("minion_id") REFERENCES "public"."minions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "treasures" table
CREATE TABLE "public"."treasures" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "name" text NULL,
  "value" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create "lair_treasures" table
CREATE TABLE "public"."lair_treasures" (
  "lair_id" text NOT NULL,
  "treasure_id" bigint NOT NULL,
  PRIMARY KEY ("lair_id", "treasure_id"),
  CONSTRAINT "fk_lair_treasures_lair" FOREIGN KEY ("lair_id") REFERENCES "public"."lairs" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_lair_treasures_treasure" FOREIGN KEY ("treasure_id") REFERENCES "public"."treasures" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
