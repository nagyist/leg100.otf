package sql

import (
	"context"

	"github.com/leg100/otf"
	"github.com/leg100/otf/sql/pggen"
)

// CreateVCSProvider inserts an agent token, associating it with an organization
func (db *DB) CreateVCSProvider(ctx context.Context, provider *otf.VCSProvider) error {
	_, err := db.InsertVCSProvider(ctx, pggen.InsertVCSProviderParams{
		VCSProviderID:    String(provider.ID()),
		Token:            String(provider.Token()),
		Name:             String(provider.Name()),
		Cloud:            String(provider.Cloud()),
		OrganizationName: String(provider.OrganizationName()),
		CreatedAt:        Timestamptz(provider.CreatedAt()),
	})
	return err
}

func (db *DB) GetVCSProvider(ctx context.Context, id string) (*otf.VCSProvider, error) {
	provider, err := db.FindVCSProvider(ctx, String(id))
	if err != nil {
		return nil, databaseError(err)
	}
	return otf.UnmarshalVCSProviderRow(otf.VCSProviderRow(provider)), nil
}

func (db *DB) ListVCSProviders(ctx context.Context, organization string) ([]*otf.VCSProvider, error) {
	rows, err := db.FindVCSProviders(ctx, String(organization))
	if err != nil {
		return nil, databaseError(err)
	}
	var unmarshalled []*otf.VCSProvider
	for _, r := range rows {
		unmarshalled = append(unmarshalled, otf.UnmarshalVCSProviderRow(otf.VCSProviderRow(r)))
	}
	return unmarshalled, nil
}

// DeleteVCSProvider deletes an agent token.
func (db *DB) DeleteVCSProvider(ctx context.Context, id string) error {
	_, err := db.DeleteVCSProviderByID(ctx, String(id))
	if err != nil {
		return databaseError(err)
	}
	return nil
}
