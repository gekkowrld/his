-- Add search functionality

SELECT
  id,
  first_name,
  last_name
FROM
  client_fts
WHERE
  client_fts MATCH ?
ORDER BY
  rank;-- Add search functionality

SELECT id, first_name, last_name
FROM client_fts
WHERE client_fts MATCH ?
ORDER BY rank;

