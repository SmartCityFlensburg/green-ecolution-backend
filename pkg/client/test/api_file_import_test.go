/*
Green Space Management API

Testing FileImportAPIService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package client

import (
	"context"
	openapiclient "github.com/green-ecolution/green-ecolution-backend/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_client_FileImportAPIService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test FileImportAPIService ImportTreesFromCsv", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		httpRes, err := apiClient.FileImportAPI.ImportTreesFromCsv(context.Background()).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
