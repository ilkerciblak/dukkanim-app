package viladition

import (
	internal "dukkanim-api/pkg/viladition/internal"
)

func String(value any) internal.IStringValidator {
	stringValidation := internal.NewStringViladator(value)
	return stringValidation
}

func Email(value any) internal.EmailViladatorInterface {
	return internal.NewEmailViladator(value)
}
