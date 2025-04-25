/**
Define the schema to be used by the program.
It is expected to be run once on startup, its safe to do so.
This will be modelled around SQLite and Go types.
	https://sqlite.org/datatype3.html
*/

CREATE TABLE IF NOT EXISTS health_program (
	-- This is supposed to be UUID 7, but we'll treat it as a string
	-- 	https://uuid7.com/
	-- The current implementation uses approximately 36 chars, 16 bytes internally
	-- 	according to google go implementation.
	-- This stores it as a string representation as it will be viewed that way by users.
	-- Alternate definition can be `id BLOB(16) PRIMARY KEY`
	id VARCHAR(40) PRIMARY KEY UNIQUE NOT NULL,

	-- Few names are larger than this! (should be UNIQUE)
	name VARCHAR(100) NOT NULL,

	-- A short text about the program
	description TEXT,

	-- Stored as seconds from UNIX epoch, for easy conversion
	start_date INTEGER NOT NULL,

	-- If NULL, then it will never end.
	end_date INTEGER
);

CREATE TABLE IF NOT EXISTS client (
	id VARCHAR(40) PRIMARY KEY UNIQUE NOT NULL,

	-- Allow the name to be as long as possible (uniocode can easily be accomodated)
	first_name TEXT NOT NULL,
	middle_name TEXT, -- Some people have only 2 names
	last_name TEXT NOT NULL,

	-- This is date of birth, since sqlite doesn't have datetime, just
	-- 	request the individual segments
	day INTEGER NOT NULL,
	month INTEGER NOT NULL,
	year INTEGER NOT NULL
);

-- Follows the guidance of:
-- 	https://www.geeksforgeeks.org/mysql-on-delete-cascade-constraint/
CREATE TABLE IF NOT EXISTS client_program (
	client_id VARCHAR(40) NOT NULL,
	program_id VARCHAR(40) NOT NULL,
	PRIMARY KEY (client_id, program_id),
	FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE,
	FOREIGN KEY (program_id) REFERENCES health_program(id) ON DELETE CASCADE
);
