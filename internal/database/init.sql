create table if not exists order_data (
    order_uid varchar(30) unique primary key,
    track_number varchar(30),
    entry varchar(10),
    locale varchar(4),
    internal_signature varchar(30),
    customer_id varchar(20),
    delivery_service varchar(20),
    shard_key varchar(10),
    sm_id int,
    date_created varchar(30),
    oof_shard varchar(10)
);

create table if not exists delivery (
    order_uid varchar(30) unique primary key,
    name varchar(20),
    phone varchar(20),
    zip varchar(10),
    city varchar(50),
    address varchar(50),
    region varchar(50),
    email varchar(50)
);

create table if not exists payment (
    order_uid varchar(30) unique primary key,
    transaction_id varchar(30),
    request_id varchar(30),
    currency varchar(5),
    provider varchar(10),
    amount float8,
    payment_dt int,
    bank varchar(20),
    delivery_cost float8,
    goods_total int,
    custom_fee float8
);

create table if not exists item (
    order_uid varchar(30),
    chrt_id int,
    track_number varchar(30),
    price float8,
    rid varchar(30),
    name varchar(30),
    sale float8,
    size varchar(10),
    total_price float8,
    nm_id int,
    brand varchar(30),
    status int
);
create table if not exists item (
    order_uid varchar(30),
    chrt_id int,
    track_number varchar(30),
    price float8,
    rid varchar(30),
    name varchar(30),
    sale float8,
    size varchar(10),
    total_price float8,
    nm_id int,
    brand varchar(30),
    status int
);