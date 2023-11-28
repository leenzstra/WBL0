CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(19) PRIMARY KEY NOT NULL,
    track_number VARCHAR NOT NULL,
    entry VARCHAR NOT NULL,
    delivery JSONB NOT NULL,
    payment JSONB NOT NULL,
    items JSONB[] NOT NULL,
    locale VARCHAR NOT NULL,
    internal_signature VARCHAR NOT NULL,
    customer_id VARCHAR NOT NULL,
    delivery_service VARCHAR NOT NULL,
    shardkey VARCHAR NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR NOT NULL
);