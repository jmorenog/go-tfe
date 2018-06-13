package tfe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariablesList(t *testing.T) {
	client := testClient(t)

	orgTest, orgTestCleanup := createOrganization(t, client)
	defer orgTestCleanup()

	wTest, _ := createWorkspace(t, client, orgTest)

	vTest1, _ := createVariable(t, client, wTest)
	vTest2, _ := createVariable(t, client, wTest)

	t.Run("with valid options", func(t *testing.T) {
		vs, err := client.Variables.List(VariableListOptions{
			Organization: String(orgTest.Name),
			Workspace:    String(wTest.Name),
		})
		require.NoError(t, err)
		assert.Contains(t, vs, vTest1)
		assert.Contains(t, vs, vTest2)
	})

	t.Run("with list options", func(t *testing.T) {
		t.Skip(" w")
		// Request a page number which is out of range. The result should
		// be successful, but return no results if the paging options are
		// properly passed along.
		vs, err := client.Variables.List(VariableListOptions{
			ListOptions: ListOptions{
				PageNumber: 999,
				PageSize:   100,
			},
			Organization: String(orgTest.Name),
			Workspace:    String(wTest.Name),
		})
		require.NoError(t, err)
		assert.Empty(t, vs)
	})

	t.Run("when options is missing an organization", func(t *testing.T) {
		vs, err := client.Variables.List(VariableListOptions{
			Workspace: String(wTest.Name),
		})
		assert.Nil(t, vs)
		assert.EqualError(t, err, "Organization is required")
	})

	t.Run("when options is missing an workspace", func(t *testing.T) {
		vs, err := client.Variables.List(VariableListOptions{
			Organization: String(orgTest.Name),
		})
		assert.Nil(t, vs)
		assert.EqualError(t, err, "Workspace is required")
	})
}

func TestVariablesCreate(t *testing.T) {
	client := testClient(t)

	wTest, wTestCleanup := createWorkspace(t, client, nil)
	defer wTestCleanup()

	t.Run("with valid options", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(randomString(t)),
			Value:     String(randomString(t)),
			Category:  Category(CategoryTerraform),
			Workspace: wTest,
		}

		v, err := client.Variables.Create(options)
		require.NoError(t, err)

		assert.NotEmpty(t, v.ID)
		assert.Equal(t, *options.Key, v.Key)
		assert.Equal(t, *options.Value, v.Value)
		assert.Equal(t, *options.Category, v.Category)
		// The workspace isn't returned correcly by the API.
		// assert.Equal(t, *options.Workspace, v.Workspace)
	})

	t.Run("when options is missing key", func(t *testing.T) {
		options := VariableCreateOptions{
			Value:     String(randomString(t)),
			Category:  Category(CategoryTerraform),
			Workspace: wTest,
		}

		_, err := client.Variables.Create(options)
		assert.EqualError(t, err, "Key is required")
	})

	t.Run("when options is missing value", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(randomString(t)),
			Category:  Category(CategoryTerraform),
			Workspace: wTest,
		}

		_, err := client.Variables.Create(options)
		assert.EqualError(t, err, "Value is required")
	})

	t.Run("when options is missing category", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:       String(randomString(t)),
			Value:     String(randomString(t)),
			Workspace: wTest,
		}

		_, err := client.Variables.Create(options)
		assert.EqualError(t, err, "Category is required")
	})

	t.Run("when options is missing workspace", func(t *testing.T) {
		options := VariableCreateOptions{
			Key:      String(randomString(t)),
			Value:    String(randomString(t)),
			Category: Category(CategoryTerraform),
		}

		_, err := client.Variables.Create(options)
		assert.EqualError(t, err, "Workspace is required")
	})
}

func TestVariablesUpdate(t *testing.T) {
	client := testClient(t)

	vTest, vTestCleanup := createVariable(t, client, nil)
	defer vTestCleanup()

	t.Run("with valid options", func(t *testing.T) {
		options := VariableUpdateOptions{
			Key:       String("newname"),
			Value:     String("newvalue"),
			HCL:       Bool(true),
			Sensitive: Bool(true),
		}

		v, err := client.Variables.Update(vTest.ID, options)
		require.NoError(t, err)

		assert.Equal(t, *options.Key, v.Key)
		assert.Equal(t, *options.HCL, v.HCL)
		assert.Equal(t, *options.Sensitive, v.Sensitive)
		assert.Empty(t, v.Value) // Because its now sensitive
	})

	t.Run("when updating a subset of values", func(t *testing.T) {
		options := VariableUpdateOptions{
			Key:      String("someothername"),
			Category: Category(CategoryTerraform),
			HCL:      Bool(false),
		}

		v, err := client.Variables.Update(vTest.ID, options)
		require.NoError(t, err)

		assert.Equal(t, *options.Key, v.Key)
		assert.Equal(t, *options.Category, v.Category)
		assert.Equal(t, *options.HCL, v.HCL)
	})

	t.Run("without any changes", func(t *testing.T) {
		vTest, vTestCleanup := createVariable(t, client, nil)
		defer vTestCleanup()

		v, err := client.Variables.Update(vTest.ID, VariableUpdateOptions{})
		require.NoError(t, err)

		assert.Equal(t, vTest, v)
	})

	t.Run("with invalid variable ID", func(t *testing.T) {
		_, err := client.Variables.Update(badIdentifier, VariableUpdateOptions{})
		assert.EqualError(t, err, "Invalid value for variable ID")
	})
}

func TestVariablesDelete(t *testing.T) {
	client := testClient(t)

	wTest, wTestCleanup := createWorkspace(t, client, nil)
	defer wTestCleanup()

	vTest, _ := createVariable(t, client, wTest)

	t.Run("with valid options", func(t *testing.T) {
		err := client.Variables.Delete(vTest.ID)
		assert.NoError(t, err)
	})

	t.Run("with non existing variable ID", func(t *testing.T) {
		err := client.Variables.Delete("nonexisting")
		assert.EqualError(t, err, "Resource not found")
	})

	t.Run("with invalid variable ID", func(t *testing.T) {
		err := client.Variables.Delete(badIdentifier)
		assert.EqualError(t, err, "Invalid value for variable ID")
	})
}
