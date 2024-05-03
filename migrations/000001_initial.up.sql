create table if not exists student (
    id integer primary key,
    name varchar not null,
    grade varchar[] not null,
    created_at timestamptz not null,
    updated_at timestamptz not null
);
