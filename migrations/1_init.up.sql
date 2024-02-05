CREATE TABLE IF NOT EXISTS settings
(
    id serial PRIMARY KEY,
    plus_operation_execution_time integer,
    minus_operation_execution_time integer,
    multiplication_operation_execution_time integer,
    division_operation_execution_time integer
);

CREATE TABLE IF NOT EXISTS operations
(
    id serial PRIMARY KEY,
    operation varchar(200) UNIQUE
);



