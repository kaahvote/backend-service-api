CREATE TABLE "voting_policies" (
  "id" INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE "voter_policies" (
  "id" integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE "candidate_policies" (
  "id" integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE "session_states" (
  "id" integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE "users" (
  "id" integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "public_id" varchar NOT NULL,
  "name" varchar,
  "email" varchar,
  "password" varchar,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);