CREATE TABLE builds (
	id ROWID,
	url TEXT UNIQUE,
	console_log TEXT,
	built_on VARCHAR(30),
	duration BIGINT,
	displayName VARCHAR(255),
	timeStamp DATETIME,
	result VARCHAR(20)
);
