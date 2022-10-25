package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/apiclient"
)

// SourceDefinitionModel describes the data source data model.
type SourceDefinitionModel struct {
	Id                              types.String                        `tfsdk:"id"`
	Name                            types.String                        `tfsdk:"name"`
	DockerRepository                types.String                        `tfsdk:"docker_repository"`
	DockerImageTag                  types.String                        `tfsdk:"docker_image_tag"`
	DocumentationUrl                types.String                        `tfsdk:"documentation_url"`
	ProtocolVersion                 types.String                        `tfsdk:"protocol_version"`
	ReleaseStage                    types.String                        `tfsdk:"release_stage"`
	ReleaseDate                     types.String                        `tfsdk:"release_date"`
	SourceType                      types.String                        `tfsdk:"source_type"`
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

func FlattenSourceDefinition(sourceDefinition *apiclient.SourceDefinition) SourceDefinitionModel {
	var data SourceDefinitionModel

	data.Id = types.String{Value: sourceDefinition.SourceDefinitionId}
	data.Name = types.String{Value: sourceDefinition.Name}
	data.DockerRepository = types.String{Value: sourceDefinition.DockerRepository}
	data.DockerImageTag = types.String{Value: sourceDefinition.DockerImageTag}
	data.DocumentationUrl = types.String{Value: sourceDefinition.DocumentationUrl}
	data.ProtocolVersion = types.String{Value: sourceDefinition.ProtocolVersion}
	data.ReleaseStage = types.String{Value: sourceDefinition.ReleaseStage}
	data.SourceType = types.String{Value: sourceDefinition.SourceType}

	if sourceDefinition.ResourceRequirements != nil {
		if reqs := sourceDefinition.ResourceRequirements.Default; reqs != nil {
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

		if reqs := *sourceDefinition.ResourceRequirements.JobSpecific; len(reqs) > 0 {
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

	return data
}
