package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// valueBasedAlsoRequiresAttributeValidator is the underlying struct implementing ValueBasedAlsoRequires.
type valueBasedAlsoRequiresAttributeValidator struct {
	value           string
	pathExpressions path.Expressions
}

// ValueBasedAlsoRequires checks that a path.Expression has a non-null value,
// if the current attribute also has a non-null value that equals the given value.
//
// This implements the validation logic declaratively within the tfsdk.Schema.
//
// Relative path.Expression will be resolved against the validated attribute.
func ValueBasedAlsoRequires(value string, attributePaths ...path.Expression) tfsdk.AttributeValidator {
	return &valueBasedAlsoRequiresAttributeValidator{value, attributePaths}
}

var _ tfsdk.AttributeValidator = (*valueBasedAlsoRequiresAttributeValidator)(nil)

func (av valueBasedAlsoRequiresAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av valueBasedAlsoRequiresAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set to this value: %s, also these are set: %q", av.value, av.pathExpressions)
}

func (av valueBasedAlsoRequiresAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	// If attribute configuration is null, there is nothing else to validate
	if req.AttributeConfig.IsNull() {
		return
	}

	// If attribute configuration equals the value, there is nothing else to validate
	if req.AttributeConfig.(types.String).Value != av.value {
		return
	}

	expressions := req.AttributePathExpression.MergeExpressions(av.pathExpressions...)

	for _, expression := range expressions {
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)

		res.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(req.AttributePath) {
				continue
			}

			var mpVal attr.Value
			diags := req.Config.GetAttribute(ctx, mp, &mpVal)
			res.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// Delay validation until all involved attribute have a known value
			if mpVal.IsUnknown() {
				return
			}

			if mpVal.IsNull() {
				res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.AttributePath,
					fmt.Sprintf("Attribute %q must be specified when %q is specified", mp, req.AttributePath),
				))
			}
		}
	}
}
