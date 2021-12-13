package main

import (
	"fmt"
	"os"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	a "github.com/microsoft/kiota/authentication/go/azure"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/me/calendarview"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "application returned an error: %v\n", err)
		os.Exit(1)
	}
}

func toInt32Ptr(n int32) *int32 {
	return &n
}

func toStrPtr(s string) *string {
	return &s
}

func run() error {
	cred, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{})
	if err != nil {
		return fmt.Errorf("error creating credentials: %w", err)
	}

	auth, err := a.NewAzureIdentityAuthenticationProvider(cred)
	if err != nil {
		return fmt.Errorf("error authentication provider: %w", err)
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		return fmt.Errorf("error creating adapter: %w", err)
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)

	startDateTime := "2021-12-12T20:57:37.046Z"
	endDateTime := "2021-12-19T20:57:37.046Z"

	query := calendarview.CalendarViewRequestBuilderGetQueryParameters{
		StartDateTime: &startDateTime,
		EndDateTime:   &endDateTime,
	}

	options := calendarview.CalendarViewRequestBuilderGetOptions{
		Q: &query,
	}

	result, err := client.Me().CalendarView().Get(&options)
	if err != nil {
		return fmt.Errorf("get req failed: %w", err)
	}

	for _, event := range result.GetValue() {
		subject := event.GetSubject()
		startTime := event.GetStart()
		endTime := event.GetEnd()
		fmt.Printf("%s (%s - %s)\n", *subject, *startTime.GetDateTime(), *endTime.GetDateTime())
	}

	return nil
}
