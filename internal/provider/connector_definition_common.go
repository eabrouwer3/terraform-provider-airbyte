package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// ConnectorDefinitionModel describes the data source data model.
type ConnectorDefinitionModel struct {
	Id                              types.String                        `tfsdk:"id"`
	Name                            types.String                        `tfsdk:"name"`
	DockerRepository                types.String                        `tfsdk:"docker_repository"`
	DockerImageTag                  types.String                        `tfsdk:"docker_image_tag"`
	DocumentationUrl                types.String                        `tfsdk:"documentation_url"`
	ProtocolVersion                 types.String                        `tfsdk:"protocol_version"`
	ReleaseStage                    types.String                        `tfsdk:"release_stage"`
	ReleaseDate                     types.String                        `tfsdk:"release_date"`
	DefaultResourceRequirements     *resourceRequirementsModel          `tfsdk:"default_resource_requirements"`
	JobSpecificResourceRequirements *[]jobSpecResourceRequirementsModel `tfsdk:"job_specific_resource_requirements"`
}

type resourceRequirementsModel struct {
	CPURequest    types.String `tfsdk:"cpu_request"`
	CPULimit      types.String `tfsdk:"cpu_limit"`
	MemoryRequest types.String `tfsdk:"memory_request"`
	MemoryLimit   types.String `tfsdk:"memory_limit"`
}

type jobSpecResourceRequirementsModel struct {
	JobType       types.String `tfsdk:"job_type"`
	CPURequest    types.String `tfsdk:"cpu_request"`
	CPULimit      types.String `tfsdk:"cpu_limit"`
	MemoryRequest types.String `tfsdk:"memory_request"`
	MemoryLimit   types.String `tfsdk:"memory_limit"`
}

func FlattenConnectorDefinition(connectorDefinition *apiclient.ConnectorDefinition) (*ConnectorDefinitionModel, error) {
	var data ConnectorDefinitionModel

	if connectorDefinition.SourceDefinitionId != "" {
		data.Id = types.String{Value: connectorDefinition.SourceDefinitionId}
	} else if connectorDefinition.DestinationDefinitionId != "" {
		data.Id = types.String{Value: connectorDefinition.DestinationDefinitionId}
	} else {
		return nil, fmt.Errorf("either SourceDefinitionId or DestinationDefinitionId must be set, not empty")
	}
	data.Name = types.String{Value: connectorDefinition.Name}
	data.DockerRepository = types.String{Value: connectorDefinition.DockerRepository}
	data.DockerImageTag = types.String{Value: connectorDefinition.DockerImageTag}
	data.DocumentationUrl = types.String{Value: connectorDefinition.DocumentationUrl}
	data.ProtocolVersion = types.String{Value: connectorDefinition.ProtocolVersion}
	data.ReleaseStage = types.String{Value: connectorDefinition.ReleaseStage}

	if connectorDefinition.ResourceRequirements != nil {
		if reqs := connectorDefinition.ResourceRequirements.Default; reqs != nil {
			req := resourceRequirementsModel{}
			if reqs.CPURequest != "" {
				req.CPURequest = types.String{Value: reqs.CPURequest}
			} else {
				req.CPURequest = types.String{Null: true}
			}
			if reqs.CPULimit != "" {
				req.CPULimit = types.String{Value: reqs.CPULimit}
			} else {
				req.CPULimit = types.String{Null: true}
			}
			if reqs.MemoryRequest != "" {
				req.MemoryRequest = types.String{Value: reqs.MemoryRequest}
			} else {
				req.MemoryRequest = types.String{Null: true}
			}
			if reqs.MemoryLimit != "" {
				req.MemoryLimit = types.String{Value: reqs.MemoryLimit}
			} else {
				req.MemoryLimit = types.String{Null: true}
			}
			data.DefaultResourceRequirements = &req
		}

		if reqs := *connectorDefinition.ResourceRequirements.JobSpecific; len(reqs) > 0 {
			var reqData []jobSpecResourceRequirementsModel
			for _, req := range reqs {
				jobReq := jobSpecResourceRequirementsModel{
					JobType: types.String{Value: req.JobType},
				}
				if req.ResourceRequirements.CPURequest != "" {
					jobReq.CPURequest = types.String{Value: req.ResourceRequirements.CPURequest}
				} else {
					jobReq.CPURequest = types.String{Null: true}
				}
				if req.ResourceRequirements.CPULimit != "" {
					jobReq.CPULimit = types.String{Value: req.ResourceRequirements.CPULimit}
				} else {
					jobReq.CPULimit = types.String{Null: true}
				}
				if req.ResourceRequirements.MemoryRequest != "" {
					jobReq.MemoryRequest = types.String{Value: req.ResourceRequirements.MemoryRequest}
				} else {
					jobReq.MemoryRequest = types.String{Null: true}
				}
				if req.ResourceRequirements.MemoryLimit != "" {
					jobReq.MemoryLimit = types.String{Value: req.ResourceRequirements.MemoryLimit}
				} else {
					jobReq.MemoryLimit = types.String{Null: true}
				}
				reqData = append(reqData, jobReq)
			}
			data.JobSpecificResourceRequirements = &reqData
		}
	}

	return &data, nil
}

func GetCommonConnectorDefinitionFields(data ConnectorDefinitionModel) apiclient.CommonConnectorDefinitionFields {
	return apiclient.CommonConnectorDefinitionFields{
		Name:                 data.Name.Value,
		DockerRepository:     data.DockerRepository.Value,
		DockerImageTag:       data.DockerImageTag.Value,
		DocumentationUrl:     data.DocumentationUrl.Value,
		ResourceRequirements: getResourceRequirementFields(data),
	}
}

func getResourceRequirementFields(data ConnectorDefinitionModel) *apiclient.ResourceRequirements {
	if data.DefaultResourceRequirements != nil || data.JobSpecificResourceRequirements != nil {
		reqBody := apiclient.ResourceRequirements{}

		if data.DefaultResourceRequirements != nil {
			reqs := apiclient.ResourceRequirementsOptions{}
			if v := data.DefaultResourceRequirements.CPURequest; !v.IsUnknown() {
				reqs.CPURequest = v.Value
			}
			if v := data.DefaultResourceRequirements.CPULimit; !v.IsUnknown() {
				reqs.CPULimit = v.Value
			}
			if v := data.DefaultResourceRequirements.MemoryRequest; !v.IsUnknown() {
				reqs.MemoryRequest = v.Value
			}
			if v := data.DefaultResourceRequirements.MemoryLimit; !v.IsUnknown() {
				reqs.MemoryLimit = v.Value
			}
			reqBody.Default = &reqs
		}

		if data.JobSpecificResourceRequirements != nil {
			var reqs []apiclient.JobSpecificResourceRequirements
			for _, req := range *data.JobSpecificResourceRequirements {
				js := apiclient.JobSpecificResourceRequirements{
					JobType: req.JobType.Value,
				}
				if !req.CPURequest.IsUnknown() {
					js.ResourceRequirements.CPURequest = req.CPURequest.Value
				}
				if !req.CPULimit.IsUnknown() {
					js.ResourceRequirements.CPULimit = req.CPULimit.Value
				}
				if !req.MemoryRequest.IsUnknown() {
					js.ResourceRequirements.MemoryRequest = req.MemoryRequest.Value
				}
				if !req.MemoryLimit.IsUnknown() {
					js.ResourceRequirements.MemoryLimit = req.MemoryLimit.Value
				}
				reqs = append(reqs, js)
			}
			reqBody.JobSpecific = &reqs
		}

		return &reqBody
	}
	return nil
}
