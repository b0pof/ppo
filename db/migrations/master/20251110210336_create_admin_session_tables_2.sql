-- +goose Up
-- +goose StatementBegin

-- Таблица для сессий
CREATE TABLE IF NOT EXISTS goadmin_session (
                                               "id" varchar(100) NOT NULL,
                                               "values" varchar(5000) NOT NULL,
                                               "created_at" timestamp NOT NULL,
                                               "updated_at" timestamp NOT NULL,
                                               PRIMARY KEY ("id")
);

-- Таблица для сайта (настройки)
CREATE TABLE IF NOT EXISTS goadmin_site (
                                            "id" integer PRIMARY KEY,
                                            "key" varchar(100) NOT NULL,
                                            "value" text NOT NULL,
                                            "description" varchar(3000) NOT NULL,
                                            "state" integer NOT NULL DEFAULT 0,
                                            "created_at" timestamp NULL,
                                            "updated_at" timestamp NULL
);

-- Таблица для пользователей
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

-- Таблица для ролей
CREATE TABLE IF NOT EXISTS goadmin_roles (
                                             "id" serial PRIMARY KEY,
                                             "name" varchar(50) NOT NULL UNIQUE,
                                             "slug" varchar(50) NOT NULL UNIQUE,
                                             "created_at" timestamp NULL,
                                             "updated_at" timestamp NULL
);

-- Таблица для разрешений
CREATE TABLE IF NOT EXISTS goadmin_permissions (
                                                   "id" serial PRIMARY KEY,
                                                   "name" varchar(50) NOT NULL UNIQUE,
                                                   "slug" varchar(50) NOT NULL UNIQUE,
                                                   "http_method" varchar(255) NULL,
                                                   "http_path" text NOT NULL,
                                                   "created_at" timestamp NULL,
                                                   "updated_at" timestamp NULL
);

-- Таблица для меню
CREATE TABLE IF NOT EXISTS goadmin_menu (
                                            "id" serial PRIMARY KEY,
                                            "parent_id" integer DEFAULT 0,
                                            "type" integer DEFAULT 0,
                                            "order" integer DEFAULT 0,
                                            "title" varchar(50) NOT NULL,
                                            "icon" varchar(50) NOT NULL,
                                            "uri" varchar(3000) NOT NULL,
                                            "header" varchar(150) NULL,
                                            "plugin_name" varchar(150) NOT NULL,
                                            "uuid" varchar(100) NULL,
                                            "created_at" timestamp NULL,
                                            "updated_at" timestamp NULL
);

-- Таблица связи пользователей и ролей
CREATE TABLE IF NOT EXISTS goadmin_role_users (
                                                  "role_id" integer NOT NULL,
                                                  "user_id" integer NOT NULL,
                                                  "created_at" timestamp NULL,
                                                  "updated_at" timestamp NULL,
                                                  PRIMARY KEY ("role_id", "user_id")
);

-- Таблица связи ролей и разрешений
CREATE TABLE IF NOT EXISTS goadmin_role_permissions (
                                                        "role_id" integer NOT NULL,
                                                        "permission_id" integer NOT NULL,
                                                        "created_at" timestamp NULL,
                                                        "updated_at" timestamp NULL,
                                                        PRIMARY KEY ("role_id", "permission_id")
);

-- Таблица связи пользователей и разрешений
CREATE TABLE IF NOT EXISTS goadmin_user_permissions (
                                                        "user_id" integer NOT NULL,
                                                        "permission_id" integer NOT NULL,
                                                        "created_at" timestamp NULL,
                                                        "updated_at" timestamp NULL,
                                                        PRIMARY KEY ("user_id", "permission_id")
);

-- Таблица для операций (логи)
CREATE TABLE IF NOT EXISTS goadmin_operation_log (
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
