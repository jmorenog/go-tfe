package tfe

import (
	"errors"
)

// Variables handles communication with the variable related methods of the
// Terraform Enterprise API.
//
// TFE API docs: https://www.terraform.io/docs/enterprise/api/variables.html
type Variables struct {
	client *Client
}

// CategoryType represents a category type.
type CategoryType string

//List all available categories.
const (
	CategoryEnv       CategoryType = "env"
	CategoryTerraform CategoryType = "terraform"
)

// Variable represents a Terraform Enterprise variable.
type Variable struct {
	ID        string       `jsonapi:"primary,vars"`
	Key       string       `jsonapi:"attr,key"`
	Value     string       `jsonapi:"attr,value"`
	Category  CategoryType `jsonapi:"attr,category"`
	HCL       bool         `jsonapi:"attr,hcl"`
	Sensitive bool         `jsonapi:"attr,sensitive"`

	// Relations
	Workspace *Workspace `jsonapi:"relation,workspace"`
}

// VariableListOptions represents the options for listing variables.
type VariableListOptions struct {
	ListOptions
	Organization *string `url:"filter[organization][name],omitempty"`
	Workspace    *string `url:"filter[workspace][name],omitempty"`
}

func (o VariableListOptions) valid() error {
	if !validString(o.Organization) {
		return errors.New("Organization is required")
	}
	if !validString(o.Workspace) {
		return errors.New("Workspace is required")
	}
	return nil
}

// List returns all variables associated with a given workspace.
func (s *Variables) List(options VariableListOptions) ([]*Variable, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("GET", "vars", &options)
	if err != nil {
		return nil, err
	}

	result, err := s.client.do(req, []*Variable{})
	if err != nil {
		return nil, err
	}

	var vs []*Variable
	for _, v := range result.([]interface{}) {
		vs = append(vs, v.(*Variable))
	}

	return vs, nil
}

// VariableCreateOptions represents the options for creating a new variable.
type VariableCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,vars"`

	// The name of the variable.
	Key *string `jsonapi:"attr,key"`

	// The value of the variable.
	Value *string `jsonapi:"attr,value"`

	// Whether this is a Terraform or environment variable.
	Category *CategoryType `jsonapi:"attr,category"`

	// Whether to evaluate the value of the variable as a string of HCL code.
	HCL *bool `jsonapi:"attr,hcl,omitempty"`

	// Whether the value is sensitive.
	Sensitive *bool `jsonapi:"attr,sensitive,omitempty"`

	// The workspace that owns the variable.
	Workspace *Workspace `jsonapi:"relation,workspace"`
}

func (o VariableCreateOptions) valid() error {
	if !validString(o.Key) {
		return errors.New("Key is required")
	}
	if !validString(o.Value) {
		return errors.New("Value is required")
	}
	if o.Category == nil {
		return errors.New("Category is required")
	}
	if o.Workspace == nil {
		return errors.New("Workspace is required")
	}
	return nil
}

// Create is used to create a new variable.
func (s *Variables) Create(options VariableCreateOptions) (*Variable, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	req, err := s.client.newRequest("POST", "vars", &options)
	if err != nil {
		return nil, err
	}

	v, err := s.client.do(req, &Variable{})
	if err != nil {
		return nil, err
	}

	return v.(*Variable), nil
}

// VariableUpdateOptions represents the options for updating a variable.
type VariableUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,vars"`

	// The name of the variable.
	Key *string `jsonapi:"attr,key,omitempty"`

	// The value of the variable.
	Value *string `jsonapi:"attr,value,omitempty"`

	// Whether this is a Terraform or environment variable.
	Category *CategoryType `jsonapi:"attr,category,omitempty"`

	// Whether to evaluate the value of the variable as a string of HCL code.
	HCL *bool `jsonapi:"attr,hcl,omitempty"`

	// Whether the value is sensitive.
	Sensitive *bool `jsonapi:"attr,sensitive,omitempty"`
}

// Update values of an existing variable.
func (s *Variables) Update(variableID string, options VariableUpdateOptions) (*Variable, error) {
	if !validStringID(&variableID) {
		return nil, errors.New("Invalid value for variable ID")
	}

	// Make sure we don't send a user provided ID.
	options.ID = variableID

	req, err := s.client.newRequest("PATCH", "vars/"+variableID, &options)
	if err != nil {
		return nil, err
	}

	v, err := s.client.do(req, &Variable{})
	if err != nil {
		return nil, err
	}

	return v.(*Variable), nil
}

// Delete a variable.
func (s *Variables) Delete(variableID string) error {
	if !validStringID(&variableID) {
		return errors.New("Invalid value for variable ID")
	}

	req, err := s.client.newRequest("DELETE", "vars/"+variableID, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)

	return err
}
