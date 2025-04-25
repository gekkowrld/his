-- Define the schema to be used by the program.
-- It is expected to be run once on startup, its safe to do so.
-- This will be modelled around SQLite and Go types.
-- 		https://sqlite.org/datatype3.html

PRAGMA auto_vacuum = FULL;

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

-- Create virtual tables that will be used during searching.
-- This are crucial for fast search
CREATE VIRTUAL TABLE IF NOT EXISTS client_fts USING FTS5(
	id UNINDEXED,
	first_name,
	middle_name,
	last_name,
	content=client,
);

-- Create triggers to mirror the actions being done on the main table.
-- This ensures that the tables are in sync so as not to produce wrong results.
CREATE TRIGGER IF NOT EXISTS client_ai AFTER INSERT ON client
BEGIN
  INSERT INTO client_fts(id, first_name, middle_name, last_name)
  VALUES (new.id, new.first_name, new.middle_name, new.last_name);
END;

CREATE TRIGGER IF NOT EXISTS client_au AFTER UPDATE ON client
BEGIN
  UPDATE client_fts
  SET first_name = new.first_name,
      middle_name = new.middle_name,
      last_name = new.last_name
  WHERE id = new.id;
END;

CREATE TRIGGER IF NOT EXISTS client_ad AFTER DELETE ON client
BEGIN
  DELETE FROM client_fts WHERE id = old.id;
END;
