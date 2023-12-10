CREATE EXTENSION "uuid-ossp";

-- ======= USERS =======

-- initialize tables
CREATE TABLE users (
    "id" uuid NOT NULL,
    "name" varchar NOT NULL,
    "email" varchar NOT NULL UNIQUE,
    "password" varchar NOT NULL,
    "phone" varchar NOT NULL UNIQUE,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- insert user
INSERT INTO "users" ("id", "name", "email", "password", "phone", "created_at", "updated_at") VALUES
    (uuid_generate_v4(), 'User Iqbal', 'iqbal@synapsis.id', '$2a$14$osYRoEiScEbr5CVTcRVNFOxsRbl0J3Z81MbYNDgXfwKXUD6.RqFDC', '082260722260', NOW(), NOW()),
    (uuid_generate_v4(), 'User Izza', 'izza@synapsis.id', '$2a$14$mAnwQu0QALnTt1YoF0gcm.sfCruND8ZFNBjRX8tBPLklxfabUgBAm','082260141517', NOW(), NOW());

-- ======= END OF USERS =======


-- ======= RESET_TOKEN =======

-- initialize tables
CREATE TABLE reset_token (
    "id" uuid NOT NULL,
    "tokens" varchar NOT NULL,
    "email" varchar NOT NULL,
    "attempt" varchar NOT NULL,
    "last_at" varchar NOT NULL UNIQUE,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- ======= END OF RESET_TOKEN =======
