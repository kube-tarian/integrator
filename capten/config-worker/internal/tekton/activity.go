package tekton

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/model"
)

var logger = logging.NewLogger()

type TektonpipelineActivty struct{}

func (c *TektonpipelineActivty) ConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	logger.Infof("Activity: %s, %s", params.Resource, params.Action)
	status, err := processConfigurationActivity(ctx, params, payload)
	if err != nil {
		return model.ResponsePayload{
			Status: status,
			Message: json.RawMessage(
				fmt.Sprintf("{\"error\": \"%s\"}", err.Error())),
		}, err
	}

	logger.Infof("tekton pipeline action %s configured", params.Action)
	return model.ResponsePayload{Status: status}, err
}

func processConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (string, error) {
	cp, err := NewTektonApp()
	if err != nil {
		return string(model.WorkFlowStatusFailed), fmt.Errorf("failed to initialize crossplane plugin")
	}

	reqLocal := &model.TektonPipelineUseCase{}
	if err := json.Unmarshal(payload, reqLocal); err != nil {
		logger.Errorf("failed to unmarshall the tekton pipeline req for %s, %v", model.TektonPipelineCreate, err)
		return string(model.WorkFlowStatusFailed), fmt.Errorf("failed to unmarshall the crossplane req for %s", model.TektonPipelineCreate)
	}

	switch params.Action {
	case model.TektonPipelineCreate:
		status, err := cp.configureProjectAndApps(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to configure tekton project for %s, %v", model.TektonPipelineCreate, err)
			return status, fmt.Errorf("failed to configure tekton project for %s", model.TektonPipelineCreate)
		}
		return status, nil
	case model.TektonPipelineSync:
		err := cp.createOrUpdateSecrets(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to update tekton project for %s, %v", model.TektonPipelineSync, err)
			return string(model.WorkFlowStatusFailed), fmt.Errorf("failed to update tekton project for %s", model.TektonPipelineSync)
		}
		return string(model.WorkFlowStatusCompleted), nil
	case model.TektonPipelineDelete:
		status, err := cp.deleteProjectAndApps(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to delete tekton project for %s, %v", model.TektonPipelineSync, err)
			return status, fmt.Errorf("failed to delete tekton project for %s", model.TektonPipelineSync)
		}
		return status, nil
	default:
		return string(model.WorkFlowStatusFailed), fmt.Errorf("invalid tekton pipeline action")
	}
}
