package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

func OrgAccountList(includeAccountIds []string) ([]types.Account, error) {

	// try to create a default config instance
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"))

	// in the event of an issue return nil output and the error; it's expected that the caller will handle
	// the error
	if err != nil {
		return nil, err
	}

	svc := organizations.NewFromConfig(cfg)

	var accountList []types.Account

	accountPage, err := svc.ListAccounts(context.TODO(), &organizations.ListAccountsInput{NextToken: nil})
	if err != nil {
		return nil, err
	}
	accountList = append(accountList, accountPage.Accounts...)

	for accountPage.NextToken != nil {
		accountPage, err = svc.ListAccounts(context.TODO(), &organizations.ListAccountsInput{NextToken: accountPage.NextToken})
		if err != nil {
			return nil, err
		}
		accountList = append(accountList, accountPage.Accounts...)
	}

	// Filter account list
	if len(includeAccountIds) > 0 {
		var filteredAccountList []types.Account
		for _, v := range accountList {
			for _, incId := range includeAccountIds {
				if *v.Id == incId {
					filteredAccountList = append(filteredAccountList, v)
				}
			}
		}
		accountList = filteredAccountList
	}

	// return the accounts list and no error
	return accountList, nil
}
