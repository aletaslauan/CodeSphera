-- name: SelectUsers :many
select * from users;

-- name: AddUser :exec
insert into users
(username, password)
values
($1, $2);
