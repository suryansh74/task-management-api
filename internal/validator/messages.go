package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// MsgForTag returns user-friendly error messages for validation tags
func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "url":
		return "Invalid URL format"
	case "uri":
		return "Invalid URI format"
	case "uuid":
		return "Invalid UUID format"
	case "uuid4":
		return "Invalid UUID v4 format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters", fe.Param())
	case "len":
		return fmt.Sprintf("Must be exactly %s characters", fe.Param())
	case "eq":
		return fmt.Sprintf("Must be equal to %s", fe.Param())
	case "ne":
		return fmt.Sprintf("Must not be equal to %s", fe.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "gt":
		return fmt.Sprintf("Must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	case "alpha":
		return "Must contain only alphabetic characters"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "numeric":
		return "Must be numeric"
	case "number":
		return "Must be a valid number"
	case "hexadecimal":
		return "Must be a valid hexadecimal"
	case "hexcolor":
		return "Must be a valid hex color"
	case "rgb":
		return "Must be a valid RGB color"
	case "rgba":
		return "Must be a valid RGBA color"
	case "hsl":
		return "Must be a valid HSL color"
	case "hsla":
		return "Must be a valid HSLA color"
	case "e164":
		return "Must be a valid E164 phone number"
	case "latitude":
		return "Must be a valid latitude"
	case "longitude":
		return "Must be a valid longitude"
	case "ssn":
		return "Must be a valid SSN"
	case "ipv4":
		return "Must be a valid IPv4 address"
	case "ipv6":
		return "Must be a valid IPv6 address"
	case "ip":
		return "Must be a valid IP address"
	case "cidr":
		return "Must be a valid CIDR notation"
	case "mac":
		return "Must be a valid MAC address"
	case "hostname":
		return "Must be a valid hostname"
	case "fqdn":
		return "Must be a valid FQDN"
	case "unique":
		return "Must contain unique values"
	case "ascii":
		return "Must contain only ASCII characters"
	case "printascii":
		return "Must contain only printable ASCII characters"
	case "multibyte":
		return "Must contain multibyte characters"
	case "datauri":
		return "Must be a valid data URI"
	case "base64":
		return "Must be valid base64"
	case "base64url":
		return "Must be valid base64 URL"
	case "isbn":
		return "Must be a valid ISBN"
	case "isbn10":
		return "Must be a valid ISBN-10"
	case "isbn13":
		return "Must be a valid ISBN-13"
	case "json":
		return "Must be valid JSON"
	case "jwt":
		return "Must be a valid JWT"
	case "html":
		return "Must not contain HTML tags"
	case "html_encoded":
		return "Must be HTML encoded"
	case "url_encoded":
		return "Must be URL encoded"
	case "lowercase":
		return "Must be lowercase"
	case "uppercase":
		return "Must be uppercase"
	case "datetime":
		return fmt.Sprintf("Must be a valid datetime in format: %s", fe.Param())
	case "timezone":
		return "Must be a valid timezone"
	case "startswith":
		return fmt.Sprintf("Must start with '%s'", fe.Param())
	case "endswith":
		return fmt.Sprintf("Must end with '%s'", fe.Param())
	case "contains":
		return fmt.Sprintf("Must contain '%s'", fe.Param())
	case "containsany":
		return fmt.Sprintf("Must contain any of: %s", fe.Param())
	case "excludes":
		return fmt.Sprintf("Must not contain '%s'", fe.Param())
	case "excludesall":
		return fmt.Sprintf("Must not contain any of: %s", fe.Param())
	default:
		return fmt.Sprintf("Field validation failed on %s", fe.Tag())
	}
}
