/**
 *  Authored: Josh Spicer <hello@joshspicer.com>
 */

package main

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type AzureIdentityExtended struct {
	*azidentity.DefaultAzureCredential
}

func (az *AzureIdentityExtended) RefreshCredential() {
	fmt.Println("Refreshing az credential")
	az.DefaultAzureCredential = InitAzure(false)
}

func InitAzure(allowFatal bool) *azidentity.DefaultAzureCredential {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		if allowFatal {
			log.Fatal(err)
		} else {
			fmt.Println(err)
		}
	}

	// TODO: Check for required DBs

	return cred
}
