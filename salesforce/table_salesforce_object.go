package salesforce

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// LIST HYDRATE FUNCTION

func listSalesforceObjectsByTable(tableName string, salesforceCols map[string]string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		client, err := connect(ctx, d)
		if err != nil {
			plugin.Logger(ctx).Error("salesforce.listSalesforceObjectsByTable", "connection error", err)
			return nil, err
		}
		if client == nil {
			plugin.Logger(ctx).Error("salesforce.listSalesforceObjectsByTable", "client_not_found: unable to generate dynamic tables because of invalid steampipe salesforce configuration", err)
			return nil, fmt.Errorf("salesforce.listSalesforceObjectsByTable: client_not_found, unable to query table %s because of invalid steampipe salesforce configuration", d.Table.Name)
		}

		requestedCols := getMatchingColumns(d.Table, d.QueryContext.Columns)
		plugin.Logger(ctx).Debug("salesforce.listSalesforceObjectsByTable", "requested_columns", requestedCols)

		query := generateQuery(requestedCols, tableName)
		condition := buildQueryFromQuals(d.Quals, d.Table.Columns, salesforceCols)

		// Push down WHERE
		if condition != "" {
			query = fmt.Sprintf("%s where %s", query, condition)
			plugin.Logger(ctx).Debug("salesforce.listSalesforceObjectsByTable", "table_name", d.Table.Name, "query_condition", condition)
		}

		// Push down ORDER BY
		if d.QueryContext.SortOrder != nil && len(d.QueryContext.SortOrder) > 0 {
			var parts []string
			for _, sc := range d.QueryContext.SortOrder {
				var order string
				switch sc.Order {
				case plugin.SortAsc:
					order = "ASC"
				case plugin.SortDesc:
					order = "DESC"
				default:
					order = ""
				}
				if order != "" {
					parts = append(parts, fmt.Sprintf("%s %s", getSalesforceColumnName(sc.Column), order))
				}
			}
			if len(parts) > 0 {
				query = fmt.Sprintf("%s order by %s", query, strings.Join(parts, ","))
			}
		}

		// Push down LIMIT
		if d.QueryContext.Limit != nil {
			limit := int32(*d.QueryContext.Limit)
			query = fmt.Sprintf("%s limit %d", query, limit)
		}

		for {
			result, err := client.Query(query)
			if err != nil {
				plugin.Logger(ctx).Error("salesforce.listSalesforceObjectsByTable", "query error", err)
				return nil, err
			}

			objectList := new([]map[string]interface{})
			err = decodeQueryResult(ctx, result.Records, objectList)
			if err != nil {
				plugin.Logger(ctx).Error("salesforce.listSalesforceObjectsByTable", "results decoding error", err)
				return nil, err
			}

			for _, object := range *objectList {
				d.StreamListItem(ctx, object)
			}

			// Paging
			if result.Done {
				break
			} else {
				query = result.NextRecordsURL
			}
		}

		return nil, nil
	}
}

//// GET HYDRATE FUNCTION

func getSalesforceObjectbyID(tableName string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		plugin.Logger(ctx).Info("salesforce.getSalesforceObjectbyID", "Table_Name", d.Table.Name)
		config := GetConfig(d.Connection)
		var id string
		if config.NamingConvention != nil && *config.NamingConvention == "api_native" {
			id = d.EqualsQualString("Id")
		} else {
			id = d.EqualsQualString("id")
		}
		if strings.TrimSpace(id) == "" {
			return nil, nil
		}

		client, err := connect(ctx, d)
		if err != nil {
			plugin.Logger(ctx).Error("salesforce.getSalesforceObjectbyID", "connection error", err)
			return nil, err
		}
		if client == nil {
			plugin.Logger(ctx).Error("salesforce.getSalesforceObjectbyID", "client_not_found: unable to generate dynamic tables because of invalid steampipe salesforce configuration", err)
			return nil, fmt.Errorf("salesforce.getSalesforceObjectbyID: client_not_found, unable to query table %s because of invalid steampipe salesforce configuration", d.Table.Name)
		}

		obj := client.SObject(tableName).Get(id)
		if obj == nil {
			// Object doesn't exist, handle the error
			plugin.Logger(ctx).Warn("salesforce.getSalesforceObjectbyID", fmt.Sprintf("%s with id \"%s\" not found", tableName, id))
			return nil, nil
		}

		object := new(map[string]interface{})
		err = decodeQueryResult(ctx, obj, object)
		if err != nil {
			plugin.Logger(ctx).Error("salesforce.getSalesforceObjectbyID", "result decoding error", err)
			return nil, err
		}

		return *object, nil
	}
}

//// TRANSFORM FUNCTION

func getFieldFromSObjectMap(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	ls := d.HydrateItem.(map[string]interface{})
	return ls[param], nil
}

func getFieldFromSObjectMapByColumnName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	salesforceColumnName := getSalesforceColumnName(d.ColumnName)
	ls := d.HydrateItem.(map[string]interface{})
	return ls[salesforceColumnName], nil
}

func getMatchingColumns(table *plugin.Table, columnNames []string) []*plugin.Column {
	var matchingColumns []*plugin.Column

	columnMap := make(map[string]*plugin.Column)
	for _, col := range table.Columns {
		columnMap[col.Name] = col
	}

	for _, colName := range columnNames {
		if col, found := columnMap[colName]; found {
			matchingColumns = append(matchingColumns, col)
		}
	}

	return matchingColumns
}
