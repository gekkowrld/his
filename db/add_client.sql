-- Add a client into sql
INSERT INTO
  client (
    id,
    first_name,
    middle_name,
    last_name,
    day,
    month,
    year
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?)
