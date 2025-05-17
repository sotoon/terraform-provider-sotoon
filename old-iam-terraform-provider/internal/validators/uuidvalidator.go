package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	uuid "github.com/satori/go.uuid"
)

type uuidValidator struct {
}

// Validator validates that an string Attribute's value is in form of uuid.
func (v *uuidValidator) Description(_ context.Context) string {
	return "value must be in form of a uuid"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v *uuidValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *uuidValidator) ValidateString(ctx context.Context, request validator.StringRequest, res *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	uuidString := request.ConfigValue.ValueString()
	_, err := uuid.FromString(uuidString)
	if err != nil {
		res.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			request.Config.Raw.String(),
			uuidString,
		))
	}
}

func UUID() validator.String {
	return &uuidValidator{}
}
