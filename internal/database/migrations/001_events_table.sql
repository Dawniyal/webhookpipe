create table events (
id uuid primary key default uuidv7(),
endpoint_id text not null,
payload JSONB not null,
status TEXT not null,
created_at timestamp not null default CURRENT_TIMESTAMP
);
