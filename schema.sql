CREATE TABLE Posts (
	id UUID PRIMARY KEY,
	title TEXT,
	subtitle TEXT,
	content TEXT,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);

CREATE TABLE Users (
	id UUID PRIMARY KEY,
	name TEXT,
	email TEXT UNIQUE,
	password BYTEA
);
