-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS goadmin_session (
                                               "id" varchar(100) NOT NULL,
                                               "values" varchar(5000) NOT NULL,
                                               "created_at" timestamp NOT NULL,
                                               "updated_at" timestamp NOT NULL,
                                               PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS goadmin_users (
                                             "id" serial PRIMARY KEY,
                                             "username" varchar(100) NOT NULL UNIQUE,
                                             "password" varchar(100) NOT NULL,
                                             "name" varchar(100) NOT NULL,
                                             "avatar" varchar(255) NULL,
                                             "remember_token" varchar(100) NULL,
                                             "created_at" timestamp NULL,
                                             "updated_at" timestamp NULL
);

CREATE TABLE IF NOT EXISTS goadmin_roles (
                                             "id" serial PRIMARY KEY,
                                             "name" varchar(50) NOT NULL UNIQUE,
                                             "slug" varchar(50) NOT NULL UNIQUE,
                                             "created_at" timestamp NULL,
                                             "updated_at" timestamp NULL
);

CREATE TABLE IF NOT EXISTS goadmin_permissions (
                                                   "id" serial PRIMARY KEY,
                                                   "name" varchar(50) NOT NULL UNIQUE,
                                                   "slug" varchar(50) NOT NULL UNIQUE,
                                                   "http_method" varchar(255) NULL,
                                                   "http_path" text NOT NULL,
                                                   "created_at" timestamp NULL,
                                                   "updated_at" timestamp NULL
);

CREATE TABLE IF NOT EXISTS goadmin_role_users (
                                                  "role_id" integer NOT NULL,
                                                  "user_id" integer NOT NULL,
                                                  "created_at" timestamp NULL,
                                                  "updated_at" timestamp NULL,
                                                  PRIMARY KEY ("role_id", "user_id")
);

CREATE TABLE IF NOT EXISTS goadmin_role_permissions (
                                                        "role_id" integer NOT NULL,
                                                        "permission_id" integer NOT NULL,
                                                        "created_at" timestamp NULL,
                                                        "updated_at" timestamp NULL,
                                                        PRIMARY KEY ("role_id", "permission_id")
);

CREATE TABLE IF NOT EXISTS goadmin_user_permissions (
                                                        "user_id" integer NOT NULL,
                                                        "permission_id" integer NOT NULL,
                                                        "created_at" timestamp NULL,
                                                        "updated_at" timestamp NULL,
                                                        PRIMARY KEY ("user_id", "permission_id")
);

CREATE TABLE IF NOT EXISTS goadmin_operations (
                                                  "id" serial PRIMARY KEY,
                                                  "user_id" integer NOT NULL,
                                                  "path" varchar(255) NOT NULL,
                                                  "method" varchar(10) NOT NULL,
                                                  "ip" varchar(120) NOT NULL,
                                                  "input" text NOT NULL,
                                                  "created_at" timestamp NULL,
                                                  "updated_at" timestamp NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
