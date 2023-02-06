package organization

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/leg100/otf"
	"github.com/leg100/otf/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganization_Create(t *testing.T) {
	ctx := context.Background()
	db := sql.NewTestDB(t)
	orgDB := newDB(db)
	org := NewTestOrganization(t)

	t.Cleanup(func() {
		orgDB.delete(ctx, org.Name())
	})

	err := db.CreateOrganization(ctx, org)
	require.NoError(t, err)

	t.Run("Duplicate", func(t *testing.T) {
		err := db.CreateOrganization(context.Background(), org)
		require.Equal(t, otf.ErrResourceAlreadyExists, err)
	})
}

func TestOrganization_Update(t *testing.T) {
	db := sql.NewTestDB(t)
	org := sql.CreateTestOrganization(t, db)

	newName := uuid.NewString()
	org, err := db.UpdateOrganization(context.Background(), org.Name(), func(org *otf.Organization) error {
		otf.UpdateOrganizationFromOpts(org, otf.OrganizationUpdateOptions{Name: &newName})
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, newName, org.Name())
}

func TestOrganization_Get(t *testing.T) {
	ctx := context.Background()
	db := sql.NewTestDB(t)
	org := sql.CreateTestOrganization(t, db)

	t.Run("by name", func(t *testing.T) {
		got, err := db.GetOrganization(ctx, org.Name())
		require.NoError(t, err)

		assert.Equal(t, org.Name(), got.Name())
		assert.Equal(t, org.ID(), got.ID())
	})

	t.Run("by id", func(t *testing.T) {
		got, err := db.GetOrganizationByID(ctx, org.ID())
		require.NoError(t, err)

		assert.Equal(t, org.Name(), got.Name())
		assert.Equal(t, org.ID(), got.ID())
	})
}

func TestOrganization_List(t *testing.T) {
	db := sql.NewTestDB(t)
	org := sql.CreateTestOrganization(t, db)

	ol, err := db.ListOrganizations(context.Background(), otf.OrganizationListOptions{})
	require.NoError(t, err)

	assert.Contains(t, ol.Items, org)
}

func TestOrganization_ListWithPagination(t *testing.T) {
	ctx := context.Background()
	db := sql.NewTestDB(t)
	_ = sql.CreateTestOrganization(t, db)
	_ = sql.CreateTestOrganization(t, db)

	t.Run("page one, two items per page", func(t *testing.T) {
		orgs, err := db.ListOrganizations(ctx, otf.OrganizationListOptions{ListOptions: otf.ListOptions{PageNumber: 1, PageSize: 2}})
		require.NoError(t, err)

		assert.Equal(t, 2, len(orgs.Items))
	})

	t.Run("page one, one item per page", func(t *testing.T) {
		orgs, err := db.ListOrganizations(ctx, otf.OrganizationListOptions{ListOptions: otf.ListOptions{PageNumber: 1, PageSize: 1}})
		require.NoError(t, err)

		assert.Equal(t, 1, len(orgs.Items))
	})

	t.Run("page two, one item per page", func(t *testing.T) {
		orgs, err := db.ListOrganizations(ctx, otf.OrganizationListOptions{ListOptions: otf.ListOptions{PageNumber: 2, PageSize: 1}})
		require.NoError(t, err)

		assert.Equal(t, 1, len(orgs.Items))
	})
}

func TestListUserOrganizations(t *testing.T) {
	ctx := context.Background()
	db := sql.NewTestDB(t)
	org1 := sql.CreateTestOrganization(t, db)
	org2 := sql.CreateTestOrganization(t, db)
	user := sql.CreateTestUser(t, db,
		otf.WithOrganizationMemberships(org1.Name(), org2.Name()))

	got, err := db.ListOrganizationsByUser(ctx, user.ID())
	require.NoError(t, err)

	assert.Contains(t, got, org1)
	assert.Contains(t, got, org2)
}

func TestOrganization_Delete(t *testing.T) {
	ctx := context.Background()
	db := sql.NewTestDB(t)
	org := sql.CreateTestOrganization(t, db)

	require.NoError(t, db.DeleteOrganization(ctx, org.Name()))

	_, err := db.GetOrganization(ctx, org.Name())
	assert.Equal(t, otf.ErrResourceNotFound, err)
}

func TestOrganization_DeleteError(t *testing.T) {
	ctx := context.Background()
	db := sql.NewTestDB(t)
	_ = sql.CreateTestOrganization(t, db)

	err := db.DeleteOrganization(ctx, "non-existent-org")

	assert.Equal(t, otf.ErrResourceNotFound, err)
}