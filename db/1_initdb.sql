create table short_link
(
    id            serial,
    original_link text unique,
    token         text unique,
    expires_at    text,
    primary key (id)
);

create index if not exists token_idx
    on short_link (token)