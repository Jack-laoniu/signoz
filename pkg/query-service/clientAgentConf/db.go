package clientAgentConf

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.signoz.io/signoz/pkg/query-service/clientAgentConf/sqlite"
	"go.signoz.io/signoz/pkg/query-service/model"
	"go.uber.org/zap"
)

// Repo handles DDL and DML ops on ingestion pipeline
type Repo struct {
	db *sqlx.DB
}

const clientConfmap = "client_confmap"

// NewRepo initiates a new ingestion repo
func NewRepo(db *sqlx.DB) Repo {
	return Repo{
		db: db,
	}
}

func (r *Repo) initDB(engine string) error {
	switch engine {
	case "sqlite3", "sqlite":
		return sqlite.InitDB(r.db)
	default:
		return fmt.Errorf("unsupported db")
	}
}

func (r *Repo) insertClientConfmap(
	ctx context.Context, postable *ClientConfmap,
) (*ClientConfmap, *model.ApiError) {
	var err error

	insertRow := &ClientConfmap{
		ID:       uuid.New().String(),
		LastConf: postable.LastConf,
	}

	insertQuery := `INSERT INTO clientconfmaps 
	(id, config_json) 
	VALUES ($1, $2)`

	_, err = r.db.ExecContext(ctx,
		insertQuery,
		insertRow.ID,
		insertRow.LastConf)

	if err != nil {
		zap.L().Error("error in inserting clientconfmap data", zap.Error(err))
		return nil, model.InternalError(errors.Wrap(err, "failed to insert clientconfmap"))
	}

	return insertRow, nil
}

// getPipelinesByVersion returns pipelines associated with a given version
func (r *Repo) getPipelinesByVersion(
	ctx context.Context, version int,
) ([]ClientConfmap, []error) {
	var errors []error
	confmaps := []ClientConfmap{}

	versionQuery := `SELECT r.id, 
		r.config_json,
		r.created_at,
		FROM clientconfmaps r,
			 agent_config_elements e,
			 agent_config_versions v
		WHERE r.id = e.element_id
		AND v.id = e.version_id
		AND e.element_type = $1
		AND v.version = $2
		ORDER BY order_id asc`

	err := r.db.SelectContext(ctx, &confmaps, versionQuery, clientConfmap, version)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get drop pipelines from db: %v", err)}
	}

	if len(confmaps) == 0 {
		return confmaps, nil
	}

	return confmaps, errors
}

// GetPipelines returns pipeline and errors (if any)
func (r *Repo) GetClientConfmaps(
	ctx context.Context, agent_id string,
) (*ClientConfmap, *model.ApiError) {
	confmap := &ClientConfmap{}

	clientconfmapsQuery := `SELECT *
		FROM clientconfmaps 
		WHERE agent_id = $1`

	err := r.db.GetContext(ctx, confmap, clientconfmapsQuery, agent_id)
	if err == sql.ErrNoRows {
		zap.L().Warn("No row found for ingestion confmap agent_id", zap.String("id", agent_id))
		return nil, model.NotFoundError(err)
	}
	if err != nil {
		zap.L().Error("failed to get ingestion confmap from db", zap.Error(err))
		return nil, model.InternalError(errors.Wrap(err, "failed to get ingestion confmap from db"))
	}
	return confmap, nil
}

func (r *Repo) UpdateConf(ctx context.Context, conf, status, msg, configId, agentId string) error {
	updateQuery := `UPDATE clientconfmaps
	set last_config = $1, 
	deploy_status = $2, 
	deploy_result = $3,
	last_hash = $4,
	WHERE agent_id = $5`

	_, err := r.db.ExecContext(ctx, updateQuery, conf, status, msg, configId, agentId)
	if err != nil {
		return model.BadRequest(err)
	}

	return nil

}

func (r *Repo) updateDeployStatusByAgentId(
	ctx context.Context, agentId string, status string, result string,
) *model.ApiError {

	updateQuery := `UPDATE clientconfmaps
	set deploy_status = $1, 
	deploy_result = $2
	WHERE agent_id = $4`

	_, err := r.db.ExecContext(ctx, updateQuery, status, result, agentId)
	if err != nil {
		zap.L().Error("failed to update deploy status", zap.Error(err))
		return model.InternalError(errors.Wrap(err, "failed to update deploy status"))
	}

	return nil
}

func (r *Repo) updateDeployStatus(ctx context.Context,
	agenId string,
	status string,
	result string,
	lastHash string,
	lastconf string) *model.ApiError {

	updateQuery := `INSERT OR REPLACE INTO clientconfmaps(
	deploy_status, 
	deploy_result,
	last_hash,
	last_config,
	agent_id
	) VALUES (
		$1,$2,$3,$4,$5
	)`

	_, err := r.db.ExecContext(ctx, updateQuery, status, result, lastHash, lastconf, agenId)
	if err != nil {
		zap.L().Error("failed to update deploy status", zap.Error(err))
		return model.BadRequest(fmt.Errorf("failed to  update deploy status"))
	}

	return nil
}
