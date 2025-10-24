-- +goose Up
-- +goose StatementBegin

create table if not exists item (
    id bigserial,
    name text not null,
    seller_id bigint,
    price int,
    rating int,
    description text default '',
    imgsrc text not null
);

create table if not exists cart_item (
    user_id bigint,
    item_id bigint,
    count int not null default 1,
    constraint unique_user_id_item_id unique (user_id, item_id)
);

create table if not exists "user" (
    id bigserial,
    role text,
    name text not null,
    login text not null unique,
    phone text,
    password text not null
);

create table if not exists "order" (
    id bigserial,
    buyer_id bigint,
    status text,
    created_at timestamptz not null default now()
);

create table if not exists order_item (
    order_id bigint,
    item_id bigint,
    count int
);

create table if not exists category (
    id bigserial,
    name text,
    parent_id bigint
);

create table if not exists item_category (
    item_id bigint,
    category_id bigint
);

create table if not exists review (
    id bigserial,
    user_id bigint,
    item_id bigint,
    rating int,
    advantages text,
    disadvantages text,
    comment text,
    created_at timestamptz not null default now()
);

CREATE OR REPLACE FUNCTION update_item_rating()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE item
    SET rating = (
        SELECT COALESCE(AVG(100 * r.rating), 0)
        FROM review r
        WHERE r.item_id = NEW.item_id
    )
    WHERE id = NEW.item_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_item_rating
    AFTER INSERT OR UPDATE OR DELETE ON review
    FOR EACH ROW
EXECUTE FUNCTION update_item_rating();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
