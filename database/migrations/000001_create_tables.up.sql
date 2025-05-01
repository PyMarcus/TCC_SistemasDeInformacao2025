CREATE TABLE "atoms" (
  "id" serial PRIMARY KEY,
  "question_id" integer,
  "question" text,
  "answer" text,
  "datasets_id" integer,
  "atom_searched" varchar,
  "atom_finded" varchar,
  "is_correct" bool,
  "failed" bool,
  "dataset_id" integer,
  "error_id" integer,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp
);

CREATE TABLE "questions" (
  "id" serial PRIMARY KEY,
  "question" text,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "errors" (
  "id" serial PRIMARY KEY,
  "definition" varchar,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "datasets" (
  "id" serial PRIMARY KEY,
  "class" text,
  "atom" varchar,
  "snippet" varchar,
  "line" varchar,
  "github_link" varchar,
  "status_code" varchar,
  "marked_by_agent_one" bool,
  "marked_by_agent_two" bool
);

COMMENT ON COLUMN "datasets"."class" IS 'Class java downloaded';

ALTER TABLE "atoms" ADD FOREIGN KEY ("datasets_id") REFERENCES "datasets" ("id");

ALTER TABLE "atoms" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id");

ALTER TABLE "atoms" ADD FOREIGN KEY ("error_id") REFERENCES "errors" ("id");
