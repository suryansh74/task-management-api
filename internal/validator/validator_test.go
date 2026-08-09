package validator

import (
	"testing"
)

type sampleStruct struct {
	Name  string `validate:"required,min=2,max=10"`
	Email string `validate:"required,email"`
}

func TestValidateStruct_Valid(t *testing.T) {
	s := sampleStruct{Name: "John", Email: "john@example.com"}
	errs := ValidateStruct(s)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateStruct_MissingRequired(t *testing.T) {
	s := sampleStruct{Name: "", Email: "john@example.com"}
	errs := ValidateStruct(s)
	if len(errs) == 0 {
		t.Fatal("expected validation errors")
	}
	if _, ok := errs["Name"]; !ok {
		t.Fatalf("expected error on Name, got %v", errs)
	}
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	s := sampleStruct{Name: "John", Email: "not-an-email"}
	errs := ValidateStruct(s)
	if len(errs) == 0 {
		t.Fatal("expected validation errors")
	}
	if _, ok := errs["Email"]; !ok {
		t.Fatalf("expected error on Email, got %v", errs)
	}
}

func TestValidateStruct_MinLength(t *testing.T) {
	s := sampleStruct{Name: "J", Email: "john@example.com"}
	errs := ValidateStruct(s)
	if len(errs) == 0 {
		t.Fatal("expected validation errors for min length")
	}
}

func TestMsgForTag(t *testing.T) {
	// Smoke test that ValidateStruct uses MsgForTag without panic
	s := sampleStruct{Name: "", Email: "bad"}
	errs := ValidateStruct(s)
	if len(errs) < 1 {
		t.Fatal("expected at least one error message")
	}
	for _, msg := range errs {
		if msg == "" {
			t.Fatal("error message should not be empty")
		}
	}
}
