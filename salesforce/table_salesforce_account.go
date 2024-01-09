package salesforce

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func SalesforceAccount(ctx context.Context, dm dynamicMap, config salesforceConfig) *plugin.Table {
	tableName := "Account"
	return &plugin.Table{
		Name:        "salesforce_account",
		Description: "Represents an individual account, which is an organization or person involved with business (such as customers, competitors, and partners).",
		List: &plugin.ListConfig{
			Hydrate:    listSalesforceObjectsByTable(tableName, dm.salesforceColumns),
			KeyColumns: dm.keyColumns,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getSalesforceObjectbyID(tableName),
			KeyColumns: plugin.SingleColumn(checkNameScheme(config, dm.cols)),
		},
		Columns: mergeTableColumns(ctx, config, dm.cols, []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier of the account in Salesforce."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the account."},
			{Name: "annual_revenue", Type: proto.ColumnType_DOUBLE, Description: "Estimated annual revenue of the account."},
			{Name: "industry", Type: proto.ColumnType_STRING, Description: "Primary business of account."},
			{Name: "owner_id", Type: proto.ColumnType_STRING, Description: "The ID of the user who currently owns this account. Default value is the user logged in to the API to perform the create."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of account, for example, Customer, Competitor, or Partner."},

			// Other columns
			{Name: "account_source", Type: proto.ColumnType_STRING, Description: "The source of the account record. For example, Advertisement, Data.com, or Trade Show."},
			{Name: "created_by_id", Type: proto.ColumnType_STRING, Description: "The id of the user who created the account."},
			{Name: "created_date", Type: proto.ColumnType_TIMESTAMP, Description: "The creation date and time of the account."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Text description of the account."},
			{Name: "is_deleted", Type: proto.ColumnType_BOOL, Description: "Indicates whether the object has been moved to the Recycle Bin (true) or not (false)."},
			{Name: "last_modified_by_id", Type: proto.ColumnType_STRING, Description: "The id of the user who last changed the contact fields, including modification date and time."},
			{Name: "last_modified_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time of last modification to account."},
			{Name: "number_of_employees", Type: proto.ColumnType_DOUBLE, Description: "Number of employees working at the company represented by this account."},
			{Name: "phone", Type: proto.ColumnType_STRING, Description: "The contact's primary phone number."},
			{Name: "website", Type: proto.ColumnType_STRING, Description: "The website of this account, for example, www.acme.com."},

			// JSON columns
			{Name: "billing_address", Type: proto.ColumnType_JSON, Description: "The billing adress of the account."},
			{Name: "shipping_address", Type: proto.ColumnType_JSON, Description: "The shipping adress of the account."},
		}),
	}
}
