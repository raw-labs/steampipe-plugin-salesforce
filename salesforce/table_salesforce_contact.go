package salesforce

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func SalesforceContact(ctx context.Context, dm dynamicMap, config salesforceConfig) *plugin.Table {
	tableName := "Contact"
	return &plugin.Table{
		Name:        "salesforce_contact",
		Description: "Represents a contact, which is a person associated with an account.",
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
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of the account that's the parent of this contact."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The full name of the contact."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "ID of the account that's the parent of this contact."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The contact's email address."},
			{Name: "owner_id", Type: proto.ColumnType_STRING, Description: "The ID of the owner of the account associated with this contact."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the contact, such as CEO or Vice President."},

			// Other columns
			{Name: "created_by_id", Type: proto.ColumnType_STRING, Description: "The Created By ID."},
			{Name: "created_date", Type: proto.ColumnType_TIMESTAMP, Description: "The Created Date.", Sort: plugin.SortAll},
			{Name: "department", Type: proto.ColumnType_STRING, Description: "The contact's department."},
			{Name: "fax", Type: proto.ColumnType_STRING, Description: "The contact's fax number."},
			{Name: "is_deleted", Type: proto.ColumnType_BOOL, Description: "Indicates whether the object has been moved to the Recycle Bin (true) or not (false)."},
			{Name: "is_email_bounced", Type: proto.ColumnType_BOOL, Description: "If bounce management is activated and an email is sent to a contact, indicates whether the email bounced (true) or not (false)."},
			{Name: "jigsaw", Type: proto.ColumnType_STRING, Description: "The Data.com Key."},
			{Name: "jigsaw_contact_id", Type: proto.ColumnType_STRING, Description: "References the company's ID in Data.com. If an account has a value in this field, it means that the account was imported from Data.com."},
			{Name: "last_activity_date", Type: proto.ColumnType_TIMESTAMP, Description: "Value is the most recent of either: Due date of the most recent event logged against the record. Due date of the most recently close task associated with the record.", Sort: plugin.SortAll},
			{Name: "last_cu_request_date", Type: proto.ColumnType_TIMESTAMP, Description: "The Last Stay-in-Touch Request Date.", Sort: plugin.SortAll},
			{Name: "last_cu_update_date", Type: proto.ColumnType_TIMESTAMP, Description: "The Last Stay-in-Touch Save Date.", Sort: plugin.SortAll},
			{Name: "last_modified_by_id", Type: proto.ColumnType_STRING, Description: "The Last Modified By ID."},
			{Name: "last_modified_date", Type: proto.ColumnType_TIMESTAMP, Description: "The Last Modified Date.", Sort: plugin.SortAll},
			{Name: "last_referenced_date", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the current user last accessed this record, a record related to this record, or a list view.", Sort: plugin.SortAll},
			{Name: "last_viewed_date", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the current user last viewed this record or list view. If this value is null, the user might have only accessed this record or list view (LastReferencedDate) but not viewed it.", Sort: plugin.SortAll},
			{Name: "master_record_id", Type: proto.ColumnType_STRING, Description: "If this record was deleted as the result of a merge, this field contains the ID of the record that remains."},
			{Name: "mobile_phone", Type: proto.ColumnType_STRING, Description: "Contact's mobile phone number."},
			{Name: "phone", Type: proto.ColumnType_STRING, Description: "Buisness telephone number for the contact."},
			{Name: "photo_url", Type: proto.ColumnType_STRING, Description: "The Photo URL."},
			{Name: "reports_to_id", Type: proto.ColumnType_STRING, Description: "The Reports To ID."},
			{Name: "system_modstamp", Type: proto.ColumnType_TIMESTAMP, Description: "The System Modstamp.", Sort: plugin.SortAll},

			// JSON columns
			{Name: "mailing_address", Type: proto.ColumnType_JSON, Description: "The Mailing Address."},
		}),
	}
}
