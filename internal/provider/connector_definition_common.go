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
	DefaultResourceRequirements     *ResourceRequirementsModel          `tfsdk:"default_resource_requirements"`
	JobSpecificResourceRequirements *[]JobSpecResourceRequirementsModel `tfsdk:"job_specific_resource_requirements"`
}

type ResourceRequirementsModel struct {
	CPURequest    types.String `tfsdk:"cpu_request"`
	CPULimit      types.String `tfsdk:"cpu_limit"`
	MemoryRequest types.String `tfsdk:"memory_request"`
	MemoryLimit   types.String `tfsdk:"memory_limit"`
}

type JobSpecResourceRequirementsModel struct {
	JobType types.String `tfsdk:"job_type"`
	// See this issue for why we can't just compose ResourceRequirementsModel into here - https://github.com/hashicorp/terraform-plugin-framework/issues/309
	CPURequest    types.String `tfsdk:"cpu_request"`
	CPULimit      types.String `tfsdk:"cpu_limit"`
	MemoryRequest types.String `tfsdk:"memory_request"`
	MemoryLimit   types.String `tfsdk:"memory_limit"`
}

func FlattenConnectorDefinition(connectorDefinition *apiclient.ConnectorDefinition) (*ConnectorDefinitionModel, error) {
	var data ConnectorDefinitionModel

	if connectorDefinition.SourceDefinitionId != "" {
		data.Id = types.StringValue(connectorDefinition.SourceDefinitionId)
	} else if connectorDefinition.DestinationDefinitionId != "" {
		data.Id = types.StringValue(connectorDefinition.DestinationDefinitionId)
	} else {
		return nil, fmt.Errorf("either SourceDefinitionId or DestinationDefinitionId must be set, not empty")
	}
	data.Name = types.StringValue(connectorDefinition.Name)
	data.DockerRepository = types.StringValue(connectorDefinition.DockerRepository)
	data.DockerImageTag = types.StringValue(connectorDefinition.DockerImageTag)
	data.DocumentationUrl = types.StringValue(connectorDefinition.DocumentationUrl)
	data.ProtocolVersion = types.StringValue(connectorDefinition.ProtocolVersion)
	data.ReleaseStage = types.StringValue(connectorDefinition.ReleaseStage)

	if connectorDefinition.ResourceRequirements != nil {
		if reqs := connectorDefinition.ResourceRequirements.Default; reqs != nil {
			req := FlattenResourceRequirements(reqs)
			data.DefaultResourceRequirements = req
		}

		if reqs := *connectorDefinition.ResourceRequirements.JobSpecific; len(reqs) > 0 {
			var reqData []JobSpecResourceRequirementsModel
			for _, req := range reqs {
				reqOptions := FlattenResourceRequirements(req.ResourceRequirements)
				jobReq := JobSpecResourceRequirementsModel{
					JobType:       types.StringValue(req.JobType),
					CPULimit:      reqOptions.CPULimit,
					CPURequest:    reqOptions.CPURequest,
					MemoryLimit:   reqOptions.MemoryLimit,
					MemoryRequest: reqOptions.MemoryRequest,
				}
				reqData = append(reqData, jobReq)
			}
			data.JobSpecificResourceRequirements = &reqData
		}
	}

	return &data, nil
}

func FlattenResourceRequirements(options *apiclient.ResourceRequirementsOptions) *ResourceRequirementsModel {
	req := ResourceRequirementsModel{}

	if options.CPURequest != "" {
		req.CPURequest = types.StringValue(options.CPURequest)
	} else {
		req.CPURequest = types.StringNull()
	}
	if options.CPULimit != "" {
		req.CPULimit = types.StringValue(options.CPULimit)
	} else {
		req.CPULimit = types.StringNull()
	}
	if options.MemoryRequest != "" {
		req.MemoryRequest = types.StringValue(options.MemoryRequest)
	} else {
		req.MemoryRequest = types.StringNull()
	}
	if options.MemoryLimit != "" {
		req.MemoryLimit = types.StringValue(options.MemoryLimit)
	} else {
		req.MemoryLimit = types.StringNull()
	}

	return &req
}

func GetCommonConnectorDefinitionFields(data ConnectorDefinitionModel) apiclient.CommonConnectorDefinitionFields {
	return apiclient.CommonConnectorDefinitionFields{
		Name:                 data.Name.ValueString(),
		DockerRepository:     data.DockerRepository.ValueString(),
		DockerImageTag:       data.DockerImageTag.ValueString(),
		DocumentationUrl:     data.DocumentationUrl.ValueString(),
		ResourceRequirements: getResourceRequirementFields(data),
	}
}

func getResourceRequirementOptions(model *ResourceRequirementsModel) *apiclient.ResourceRequirementsOptions {
	if model != nil {
		reqs := apiclient.ResourceRequirementsOptions{}
		if v := model.CPURequest; !v.IsUnknown() {
			reqs.CPURequest = v.ValueString()
		}
		if v := model.CPULimit; !v.IsUnknown() {
			reqs.CPULimit = v.ValueString()
		}
		if v := model.MemoryRequest; !v.IsUnknown() {
			reqs.MemoryRequest = v.ValueString()
		}
		if v := model.MemoryLimit; !v.IsUnknown() {
			reqs.MemoryLimit = v.ValueString()
		}
		return &reqs
	}
	return nil
}

func getResourceRequirementFields(data ConnectorDefinitionModel) *apiclient.ResourceRequirements {
	if data.DefaultResourceRequirements != nil || data.JobSpecificResourceRequirements != nil {
		reqBody := apiclient.ResourceRequirements{}

		reqBody.Default = getResourceRequirementOptions(data.DefaultResourceRequirements)

		if data.JobSpecificResourceRequirements != nil {
			var reqs []apiclient.JobSpecificResourceRequirements
			for _, req := range *data.JobSpecificResourceRequirements {
				js := apiclient.JobSpecificResourceRequirements{
					JobType: req.JobType.ValueString(),
					ResourceRequirements: getResourceRequirementOptions(&ResourceRequirementsModel{
						CPURequest:    req.CPURequest,
						CPULimit:      req.CPULimit,
						MemoryRequest: req.MemoryRequest,
						MemoryLimit:   req.MemoryLimit,
					}),
				}
				reqs = append(reqs, js)
			}
			reqBody.JobSpecific = &reqs
		}

		return &reqBody
	}
	return nil
}
