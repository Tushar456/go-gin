package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CustomerValidateErrorMessage(ctx *gin.Context, validateError validator.ValidationErrors) {
	out := make(map[string]string)
	for _, ve := range validateError {
		out[ve.Field()] = customeMessageForTag(ve.Tag())
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"error": out})
}

func customeMessageForTag(tag string) string {

	switch tag {
	case "required":
		return "This field is required."
	case "email":
		return "Invalid email."
	case "lte":
		return "less than 130."
	case "gte":
		return "greater than 1."
	case "datetime":
		return "enter date in yyyy-mm-dd format."
	case "e164":
		return "phone number is invalid."
	case "alphanum":
		return "field should contains only alhabets and number"
	case "number":
		return "field should contains only number"
	}
	return ""
}
