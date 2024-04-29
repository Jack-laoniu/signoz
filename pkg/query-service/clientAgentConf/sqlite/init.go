package sqlite

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func InitDB(db *sqlx.DB) error {
	var err error
	if db == nil {
		return fmt.Errorf("invalid db connection")
	}

	table_schema := `CREATE TABLE IF NOT EXISTS clientconfmaps(
		id TEXT PRIMARY KEY,
		agent_id TEXT,
		deploy_status VARCHAR(80) NOT NULL DEFAULT 'DIRTY',
		deploy_result TEXT,
		last_hash TEXT,
		last_config TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(table_schema)
	if err != nil {
		return errors.Wrap(err, "Error in creating clientconfmaps table")
	}
	return nil
}
