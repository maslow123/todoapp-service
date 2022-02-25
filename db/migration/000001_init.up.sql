CREATE TABLE "users" (
  "id" SERIAL NOT NULL,
  "name" varchar(255) NOT NULL,
  "address" text NOT NULL,
  "pic" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(now()),
  "updated_at" timestamp NOT NULL DEFAULT('0001-01-01 00:00:00Z')
);

CREATE TABLE "accounts" (
  "id" SERIAL PRIMARY KEY,
  "email" varchar(50) NOT NULL,
  "hashed_password" varchar(150) NOT NULL,
  "user_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(now()),
  "updated_at" timestamp NOT NULL DEFAULT('0001-01-01 00:00:00Z')
);

CREATE TABLE "todos" (
  "id" SERIAL PRIMARY KEY,
  "category_id" int NOT NULL,
  "user_id" int NOT NULL,
  "title" varchar(200) NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(now()),
  "updated_at" timestamp NOT NULL DEFAULT('0001-01-01 00:00:00Z')
);

CREATE TABLE "categories" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(150) NOT NULL,
  "color" varchar(10) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(now()),
  "updated_at" timestamp NOT NULL DEFAULT('0001-01-01 00:00:00Z')
);

-- ALTER TABLE "users" ADD FOREIGN KEY ("id") REFERENCES "accounts" ("user_id");

-- ALTER TABLE "users" ADD FOREIGN KEY ("id") REFERENCES "todos" ("user_id");

-- ALTER TABLE "categories" ADD FOREIGN KEY ("id") REFERENCES "todos" ("category_id");
