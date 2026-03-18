package storage

const (
	RECIPES = `
		CREATE TABLE IF NOT EXISTS recipes (
			id INTEGER PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			energy INTEGER NOT NULL,
			hunger NUMERIC NOT NULL,
			name TEXT NOT NULL,
			resource TEXT NOT NULL,
			hash TEXT NOT NULL UNIQUE,
			str1 NUMERIC,
			str2 NUMERIC,
			agi1 NUMERIC,
			agi2 NUMERIC,
			int1 NUMERIC,
			int2 NUMERIC,
			con1 NUMERIC,
			con2 NUMERIC,
			prc1 NUMERIC,
			prc2 NUMERIC,
			csm1 NUMERIC,
			csm2 NUMERIC,
			dex1 NUMERIC,
			dex2 NUMERIC,
			wil1 NUMERIC,
			wil2 NUMERIC,
			psy1 NUMERIC,
			psy2 NUMERIC
		);

		CREATE INDEX IF NOT EXISTS recipes_name ON recipes(name);
		CREATE INDEX IF NOT EXISTS recipes_hash ON recipes(hash);
	`

	INGREDIENTS = `
		CREATE TABLE IF NOT EXISTS ingredients (
			id INTEGER PRIMARY KEY,
			recipe INTEGER NOT NULL REFERENCES recipes(id),
			name TEXT NOT NULL,
			rate INTEGER NOT NULL
		);
	`
)
