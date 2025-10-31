package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 1 {
		panic("no server provided")
	}

	configServer := os.Args[1]
	fmt.Printf("Using: %s\n", configServer)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	// Establish a connection to the Key Vault client
	client, err := azappconfig.NewClient(configServer, cred, nil)
	if err != nil {
		panic(err)
	}

	revisionsPager := client.NewListSettingsPager(
		azappconfig.SettingSelector{
			KeyFilter:   to.Ptr("*"),
			LabelFilter: to.Ptr("*"),
			Fields:      azappconfig.AllSettingFields(),
		},
		nil,
	)

	settings := []azappconfig.Setting{}

	for revisionsPager.More() {
		fmt.Println("page")
		resp, err := revisionsPager.NextPage(context.Background())
		if err != nil {
			panic(errors.Wrap(err, "failed to get paged secrets"))
		}

		settings = append(settings, resp.Settings...)
	}

	for _, setting := range settings {
		fmt.Printf("%s\n", *setting.Key)
	}

	// resp, err := client.GetSetting(
	// 	context.TODO(),
	// 	"carto/uksouth/alexanderson-3871251af/ionaDynamic.chartValues",
	// 	&azappconfig.GetSettingOptions{
	// 		Label: to.Ptr("alexanderson-3871251af"),
	// 	})

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(*resp.Key)
}
