create table credit_cards
(
    id            bigserial    not null,
    customer_name varchar(100) not null,
    card_number   varchar(16)  not null,
    cvv           varchar(6)   not null,
    expiry_month  numeric      not null,
    expiry_year   numeric      not null,
    primary key (id),
    CONSTRAINT cc_card_num_unique UNIQUE (card_number)
);