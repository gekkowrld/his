-- Get user profile, including program IDs and names they are enrolled in

SELECT
  c.id AS client_id,
  c.first_name,
  c.last_name,
  GROUP_CONCAT(cp.program_id) AS program_ids,
  GROUP_CONCAT(hp.name) AS program_names
FROM
  client c
  LEFT JOIN client_program cp ON c.id = cp.client_id
  LEFT JOIN health_program hp ON cp.program_id = hp.id
WHERE
  c.id = ?
GROUP BY
  c.id;
