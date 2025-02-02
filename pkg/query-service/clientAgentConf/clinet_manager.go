package clientAgentConf

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.signoz.io/signoz/pkg/query-service/model"
	"gopkg.in/yaml.v2"
)

var m *Manager

func init() {
	m = &Manager{}
}

type AgentFeatureType string

type Manager struct {
	Repo
	// lock to make sure only one update is sent to remote agents at a time
	lock uint32

	// For AgentConfigProvider implementation
	configSubscribers     map[string]func()
	configSubscribersLock sync.Mutex
}

type ManagerOptions struct {
	DB       *sqlx.DB
	DBEngine string
}

func Initiate(options *ManagerOptions) (*Manager, error) {
	m = &Manager{
		Repo:              Repo{options.DB},
		configSubscribers: map[string]func(){},
	}

	err := m.initDB(options.DBEngine)
	if err != nil {
		return nil, errors.Wrap(err, "could not init agentConf db")
	}
	return m, nil
}

// Report deployment status for config recommendations generated by RecommendAgentConfig
// Implements opamp.AgentConfigProvider
func (m *Manager) ReportConfigDeploymentStatus(
	agentId string,
	configId string,
	err error,
) {
	newStatus := string(Deployed)
	message := "Deployment was successful"
	if err != nil {
		newStatus = string(DeployFailed)
		message = fmt.Sprintf("%s: %s", agentId, err.Error())
	}
	m.updateDeployStatusByAgentId(
		context.Background(), agentId, newStatus, message,
	)
}

// Implements opamp.AgentConfigProvider
func (m *Manager) RecommendAgentConfig(agentId string, currentConfYaml []byte) (
	recommendedConfYaml []byte,
	// Opaque id of the recommended config, used for reporting deployment status updates
	configId string,
	err error,
) {
	var recommendation []byte
	updatedConf, apiErr := m.GetClientConfmaps(context.Background(), agentId)
	if apiErr != nil && apiErr.Typ != model.ErrorNotFound {
		return nil, "", errors.Wrap(apiErr.ToError(), fmt.Sprintf(
			"failed to generate agent config recommendation for ClientConfmaps",
		))
	}
	recommendation = currentConfYaml

	if apiErr.Typ != model.ErrorNotFound {
		updatedConf, err := yaml.Marshal(updatedConf.LastConf)
		if err != nil {
			return nil, "", errors.Wrap(apiErr.ToError(), fmt.Sprintf(
				"yaml.Marshal err",
			))
		}
		recommendation = updatedConf
	}

	// It is possible for a feature to recommend collector config
	// before any user created config versions exist.
	//
	// For example, log pipeline config for installed integrations will
	// have to be recommended even if the user hasn't created any pipelines yet

	// Do not return an empty configId even if no recommendations were made
	hash := sha256.New()
	hash.Write(recommendation)
	configId = string(hash.Sum(nil))

	m.updateDeployStatus(
		context.Background(),
		agentId,
		string(DeployInitiated),
		"Deployment has started",
		configId,
		string(recommendation),
	)
	fmt.Println(string(recommendation))
	return recommendation, configId, nil
}

// Implements opamp.AgentConfigProvider
func (m *Manager) SubscribeToConfigUpdates(callback func()) (unsubscribe func()) {
	m.configSubscribersLock.Lock()
	defer m.configSubscribersLock.Unlock()

	subscriberId := uuid.NewString()
	m.configSubscribers[subscriberId] = callback

	return func() {
		delete(m.configSubscribers, subscriberId)
	}
}

func (m *Manager) notifyConfigUpdateSubscribers() {
	m.configSubscribersLock.Lock()
	defer m.configSubscribersLock.Unlock()
	for _, handler := range m.configSubscribers {
		handler()
	}
}


func UpdateConfmap(
	ctx context.Context, conf, agentId string,
) *model.ApiError {
	hash := sha256.New()
	hash.Write([]byte(conf))
	configId := string(hash.Sum(nil))
	// we need update lastconf
	err := m.UpdateConf(ctx, conf, string(DeployInitiated),
		"Deployment has started", configId, agentId)
	if err != nil {
		return model.InternalError(errors.Wrap(err, "failed to update confmap from db"))
	}

	m.notifyConfigUpdateSubscribers()

	return nil
}
