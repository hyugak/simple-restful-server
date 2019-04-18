CREATE TABLE users
(
    id SERIAL primary key,
    name TEXT not null,
    email TEXT not null,
    created_at TIMESTAMP not null default current_timestamp,
    updated_at TIMESTAMP not null default current_timestamp
);

CREATE function set_update_time() returns opaque as '
begin
    new.updated_at := ''now'';
    return new;
end;
'language 'plpgsql';

CREATE trigger update_tri before update on users for each row execute procedure set_update_time();
