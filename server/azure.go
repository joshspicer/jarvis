/**
 *  Authored: Josh Spicer <hello@joshspicer.com>
 */

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type AzureExtended struct {
	credential *azidentity.DefaultAzureCredential
	// containerClient *azcosmos.ContainerClient
}

func (az *AzureExtended) RefreshCredential(allowFatal bool) {
	fmt.Println("Refreshing az credential")
	initAzure(allowFatal)
}

func ConnectToCosmosContainer(allowFatal bool) *azcosmos.ContainerClient {

	connectionString := os.Getenv("AZURE_COSMOS_CONNECTION_STRING")
	dbUrl := os.Getenv("AZURE_COSMOS_DB_URL")

	if connectionString == "" || dbUrl == "" {
		errMsg := "Missing AZURE_COSMOS_CONNECTION_STRING or AZURE_COSMOS_DB_URL"
		if allowFatal {
			log.Fatal(errMsg)
		} else {
			fmt.Println(errMsg)
		}
	}
	cred, _ := azcosmos.NewKeyCredential(connectionString)

	client, err := azcosmos.NewClientWithKey(dbUrl, cred, nil)
	if err != nil {
		if allowFatal {
			log.Fatal(err)
		} else {
			fmt.Println(err)
		}
	}

	container, err := client.NewContainer("primaryDB", "primaryCollection")
	if err != nil {
		if allowFatal {
			log.Fatal(err)
		} else {
			fmt.Println(err)
		}
	}

	return container
}

func initAzure(allowFatal bool) *AzureExtended {
	opts := azidentity.DefaultAzureCredentialOptions{}
	cred, err := azidentity.NewDefaultAzureCredential(&opts)
	if err != nil {
		if allowFatal {
			log.Fatal(err)
		} else {
			fmt.Println(err)
		}
	}

	extended := AzureExtended{
		credential: cred,
		// containerClient: ConnectToCosmosContainer(allowFatal),
	}

	return &extended
}
