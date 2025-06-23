-- CreateEnum
CREATE TYPE account_status AS ENUM ('active', 'inactive', 'suspended', 'pending');

-- CreateTable
CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY,
    "nick_name" TEXT NOT NULL,
    "first_name" TEXT NOT NULL,
    "last_name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "status" account_status NOT NULL DEFAULT 'pending',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "password" TEXT NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CreateIndex
CREATE UNIQUE INDEX "users_nick_name_key" ON "users"("nick_name");

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

-- Add comments to table and columns
COMMENT ON TABLE "users" IS 'Stores user account information for the financial application';
COMMENT ON COLUMN "users"."id" IS 'Unique identifier for the user';
COMMENT ON COLUMN "users"."nick_name" IS 'User''s chosen display name';
COMMENT ON COLUMN "users"."first_name" IS 'User''s first name';
COMMENT ON COLUMN "users"."last_name" IS 'User''s last name';
COMMENT ON COLUMN "users"."email" IS 'User''s email address (must be unique)';
COMMENT ON COLUMN "users"."status" IS 'Current state of the user''s account';
COMMENT ON COLUMN "users"."created_at" IS 'Timestamp when the user account was created';
COMMENT ON COLUMN "users"."password" IS 'Hashed password (never stored in plain text)';
COMMENT ON COLUMN "users"."updated_at" IS 'Timestamp when the user was last updated';