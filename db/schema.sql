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
	end_date INTEGER,

	-- This is the day that the program was created, should be GREATER
	-- than the start_date and end_date
	created_on INTEGER NOT NULL
);

