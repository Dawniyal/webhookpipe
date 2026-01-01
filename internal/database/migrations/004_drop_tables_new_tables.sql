drop table if exists events;
drop table if exists endpoints;

create type event_status as enum (
  'received',
  'pending',
  'delivered',
  'failed'
);

create table endpoint (
  id text primary key,
  target_url text not null,
  active boolean not null default true,
  created_at timestamptz not null default now()
);

create index idx_endpoint_target_url
  on endpoint(target_url);

create table event (
  id uuid primary key,
  endpoint_id text not null,
  payload jsonb not null,
  status event_status not null default 'received',
  active boolean not null default true,
  created_at timestamptz not null default now(),
  constraint event_endpoint_id_fkey
    foreign key (endpoint_id)
    references endpoint(id)
    on delete restrict
    on update cascade
);

create index idx_event_endpoint_id
  on event(endpoint_id);
