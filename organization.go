package otf

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/leg100/otf/http/jsonapi"
)

const (
	DefaultSessionTimeout    = 20160
	DefaultSessionExpiration = 20160
)

type (
	// Organization is an OTF Organization, comprising workspaces, users, etc.
	Organization struct {
		ID              string
		CreatedAt       time.Time
		UpdatedAt       time.Time
		Name            string
		SessionRemember int
		SessionTimeout  int
	}

	// OrganizationList represents a list of Organizations.
	OrganizationList struct {
		*Pagination
		Items []*Organization
	}

	OrganizationService interface {
		CreateOrganization(ctx context.Context, opts OrganizationCreateOptions) (Organization, error)
		EnsureCreatedOrganization(ctx context.Context, opts OrganizationCreateOptions) (Organization, error)
		GetOrganization(ctx context.Context, name string) (Organization, error)
		GetOrganizationJSONAPI(ctx context.Context, name string) (*jsonapi.Organization, error)
	}

	// OrganizationCreateOptions represents the options for creating an
	// organization. See dto.OrganizationCreateOptions for more details.
	OrganizationCreateOptions struct {
		Name            *string `schema:"name,required"`
		SessionRemember *int
		SessionTimeout  *int
	}

	// ListOptions represents the options for listing organizations.
	OrganizationListOptions struct {
		ListOptions
	}

	// UpdateOptions represents the options for updating an
	// organization.
	OrganizationUpdateOptions struct {
		Name            *string
		SessionRemember *int
		SessionTimeout  *int
	}
)

func NewOrganization(opts OrganizationCreateOptions) (*Organization, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	org := Organization{
		Name:            *opts.Name,
		CreatedAt:       CurrentTimestamp(),
		UpdatedAt:       CurrentTimestamp(),
		ID:              NewID("org"),
		SessionTimeout:  DefaultSessionTimeout,
		SessionRemember: DefaultSessionExpiration,
	}
	if opts.SessionTimeout != nil {
		org.SessionTimeout = *opts.SessionTimeout
	}
	if opts.SessionRemember != nil {
		org.SessionRemember = *opts.SessionRemember
	}
	return &org, nil
}

func (org *Organization) String() string { return org.ID }

func (org *Organization) Update(opts OrganizationUpdateOptions) error {
	if opts.Name != nil {
		org.Name = *opts.Name
	}
	if opts.SessionTimeout != nil {
		org.SessionTimeout = *opts.SessionTimeout
	}
	if opts.SessionRemember != nil {
		org.SessionRemember = *opts.SessionRemember
	}
	return nil
}

func (opts *OrganizationCreateOptions) Validate() error {
	if opts.Name == nil {
		return errors.New("name required")
	}
	if *opts.Name == "" {
		return errors.New("name cannot be empty")
	}
	if !ValidStringID(opts.Name) {
		return fmt.Errorf("invalid name: %s", *opts.Name)
	}
	return nil
}
