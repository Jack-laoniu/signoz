// package clientconfmap

// import (
// 	"context"

// 	"go.signoz.io/signoz/pkg/query-service/agentConf"
// 	"go.signoz.io/signoz/pkg/query-service/model"
// )

// // Controller takes care of deployment cycle of log parsing pipelines.
// type ClientConfmapController struct {

// }

// // Implements agentConf.AgentFeature interface.
// func (pc *ClientConfmapController) AgentFeatureType() agentConf.AgentFeatureType {
// 	return ClientConfmapFeatureType
// }

// 	// Recommend config for an agent based on its `currentConfYaml` and
// 	// `configVersion` for the feature's settings
// // Implements agentConf.AgentFeature interface.
// func (pc *ClientConfmapController) RecommendAgentConfig(
// 	currentConfYaml []byte,
// 	configVersion *agentConf.ConfigVersion,
// ) (
// 	recommendedConfYaml []byte,
// 	serializedSettingsUsed string,
// 	apiErr *model.ApiError,
// ) {
// 	pipelinesVersion := -1
// 	if configVersion != nil {
// 		pipelinesVersion = configVersion.Version
// 	}

// 	pipelinesResp, apiErr := pc.GetPipelinesByVersion(
// 		context.Background(), pipelinesVersion,
// 	)
// 	if apiErr != nil {
// 		return nil, "", apiErr
// 	}

// 	updatedConf, apiErr := GenerateCollectorConfigWithPipelines(
// 		currentConfYaml, pipelinesResp.Pipelines,
// 	)
// 	if apiErr != nil {
// 		return nil, "", model.WrapApiError(apiErr, "could not marshal yaml for updated conf")
// 	}

// 	rawPipelineData, err := json.Marshal(pipelinesResp.Pipelines)
// 	if err != nil {
// 		return nil, "", model.BadRequest(errors.Wrap(err, "could not serialize pipelines to JSON"))
// 	}
// //
// 	return updatedConf, string(rawPipelineData), nil
// }