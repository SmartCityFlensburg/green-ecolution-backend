/*
Green Space Management API

Testing TreeClusterAPIService

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

func Test_client_TreeClusterAPIService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test TreeClusterAPIService CreateTreeCluster", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		resp, httpRes, err := apiClient.TreeClusterAPI.CreateTreeCluster(context.Background()).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test TreeClusterAPIService DeleteTreeCluster", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var clusterId int32

		httpRes, err := apiClient.TreeClusterAPI.DeleteTreeCluster(context.Background(), clusterId).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test TreeClusterAPIService GetAllTreeClusters", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		resp, httpRes, err := apiClient.TreeClusterAPI.GetAllTreeClusters(context.Background()).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test TreeClusterAPIService GetTreeClusterById", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var clusterId int32

		resp, httpRes, err := apiClient.TreeClusterAPI.GetTreeClusterById(context.Background(), clusterId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test TreeClusterAPIService UpdateTreeCluster", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var clusterId int32

		resp, httpRes, err := apiClient.TreeClusterAPI.UpdateTreeCluster(context.Background(), clusterId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
