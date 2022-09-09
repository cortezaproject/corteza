package rdbms

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	automationType "github.com/cortezaproject/corteza-server/automation/types"
	composeType "github.com/cortezaproject/corteza-server/compose/types"
	federationType "github.com/cortezaproject/corteza-server/federation/types"
	actionlogType "github.com/cortezaproject/corteza-server/pkg/actionlog"
	discoveryType "github.com/cortezaproject/corteza-server/pkg/discovery/types"
	flagType "github.com/cortezaproject/corteza-server/pkg/flag/types"
	labelsType "github.com/cortezaproject/corteza-server/pkg/label/types"
	rbacType "github.com/cortezaproject/corteza-server/pkg/rbac"
	systemType "github.com/cortezaproject/corteza-server/system/types"
	"github.com/doug-martin/goqu/v9"
)

var (
	// actionlogTable represents actionlogs store table
	//
	// This value is auto-generated
	actionlogTable = goqu.T("actionlog")

	// actionlogSelectQuery assembles select query for fetching actionlogs
	//
	// This function is auto-generated
	actionlogSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"ts",
			"actor_ip_addr",
			"actor_id",
			"request_origin",
			"request_id",
			"resource",
			"action",
			"error",
			"severity",
			"description",
			"meta",
		).From(actionlogTable)
	}

	// actionlogInsertQuery assembles query inserting actionlogs
	//
	// This function is auto-generated
	actionlogInsertQuery = func(d goqu.DialectWrapper, res *actionlogType.Action) *goqu.InsertDataset {
		return d.Insert(actionlogTable).
			Rows(goqu.Record{
				"id":             res.ID,
				"ts":             res.Timestamp,
				"actor_ip_addr":  res.ActorIPAddr,
				"actor_id":       res.ActorID,
				"request_origin": res.RequestOrigin,
				"request_id":     res.RequestID,
				"resource":       res.Resource,
				"action":         res.Action,
				"error":          res.Error,
				"severity":       res.Severity,
				"description":    res.Description,
				"meta":           res.Meta,
			})
	}

	// actionlogUpsertQuery assembles (insert+on-conflict) query for replacing actionlogs
	//
	// This function is auto-generated
	actionlogUpsertQuery = func(d goqu.DialectWrapper, res *actionlogType.Action) *goqu.InsertDataset {
		var target = `,id`

		return actionlogInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"ts":             res.Timestamp,
						"actor_ip_addr":  res.ActorIPAddr,
						"actor_id":       res.ActorID,
						"request_origin": res.RequestOrigin,
						"request_id":     res.RequestID,
						"resource":       res.Resource,
						"action":         res.Action,
						"error":          res.Error,
						"severity":       res.Severity,
						"description":    res.Description,
						"meta":           res.Meta,
					},
				),
			)
	}

	// actionlogUpdateQuery assembles query for updating actionlogs
	//
	// This function is auto-generated
	actionlogUpdateQuery = func(d goqu.DialectWrapper, res *actionlogType.Action) *goqu.UpdateDataset {
		return d.Update(actionlogTable).
			Set(goqu.Record{
				"ts":             res.Timestamp,
				"actor_ip_addr":  res.ActorIPAddr,
				"actor_id":       res.ActorID,
				"request_origin": res.RequestOrigin,
				"request_id":     res.RequestID,
				"resource":       res.Resource,
				"action":         res.Action,
				"error":          res.Error,
				"severity":       res.Severity,
				"description":    res.Description,
				"meta":           res.Meta,
			}).
			Where(actionlogPrimaryKeys(res))
	}

	// actionlogDeleteQuery assembles delete query for removing actionlogs
	//
	// This function is auto-generated
	actionlogDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(actionlogTable).Where(ee...)
	}

	// actionlogDeleteQuery assembles delete query for removing actionlogs
	//
	// This function is auto-generated
	actionlogTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(actionlogTable)
	}

	// actionlogPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	actionlogPrimaryKeys = func(res *actionlogType.Action) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// apigwFilterTable represents apigwFilters store table
	//
	// This value is auto-generated
	apigwFilterTable = goqu.T("apigw_filters")

	// apigwFilterSelectQuery assembles select query for fetching apigwFilters
	//
	// This function is auto-generated
	apigwFilterSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_route",
			"weight",
			"kind",
			"ref",
			"enabled",
			"params",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(apigwFilterTable)
	}

	// apigwFilterInsertQuery assembles query inserting apigwFilters
	//
	// This function is auto-generated
	apigwFilterInsertQuery = func(d goqu.DialectWrapper, res *systemType.ApigwFilter) *goqu.InsertDataset {
		return d.Insert(apigwFilterTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"rel_route":  res.Route,
				"weight":     res.Weight,
				"kind":       res.Kind,
				"ref":        res.Ref,
				"enabled":    res.Enabled,
				"params":     res.Params,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			})
	}

	// apigwFilterUpsertQuery assembles (insert+on-conflict) query for replacing apigwFilters
	//
	// This function is auto-generated
	apigwFilterUpsertQuery = func(d goqu.DialectWrapper, res *systemType.ApigwFilter) *goqu.InsertDataset {
		var target = `,id`

		return apigwFilterInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_route":  res.Route,
						"weight":     res.Weight,
						"kind":       res.Kind,
						"ref":        res.Ref,
						"enabled":    res.Enabled,
						"params":     res.Params,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
						"created_by": res.CreatedBy,
						"updated_by": res.UpdatedBy,
						"deleted_by": res.DeletedBy,
					},
				),
			)
	}

	// apigwFilterUpdateQuery assembles query for updating apigwFilters
	//
	// This function is auto-generated
	apigwFilterUpdateQuery = func(d goqu.DialectWrapper, res *systemType.ApigwFilter) *goqu.UpdateDataset {
		return d.Update(apigwFilterTable).
			Set(goqu.Record{
				"rel_route":  res.Route,
				"weight":     res.Weight,
				"kind":       res.Kind,
				"ref":        res.Ref,
				"enabled":    res.Enabled,
				"params":     res.Params,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			}).
			Where(apigwFilterPrimaryKeys(res))
	}

	// apigwFilterDeleteQuery assembles delete query for removing apigwFilters
	//
	// This function is auto-generated
	apigwFilterDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(apigwFilterTable).Where(ee...)
	}

	// apigwFilterDeleteQuery assembles delete query for removing apigwFilters
	//
	// This function is auto-generated
	apigwFilterTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(apigwFilterTable)
	}

	// apigwFilterPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	apigwFilterPrimaryKeys = func(res *systemType.ApigwFilter) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// apigwRouteTable represents apigwRoutes store table
	//
	// This value is auto-generated
	apigwRouteTable = goqu.T("apigw_routes")

	// apigwRouteSelectQuery assembles select query for fetching apigwRoutes
	//
	// This function is auto-generated
	apigwRouteSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"endpoint",
			"method",
			"enabled",
			"meta",
			"rel_group",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(apigwRouteTable)
	}

	// apigwRouteInsertQuery assembles query inserting apigwRoutes
	//
	// This function is auto-generated
	apigwRouteInsertQuery = func(d goqu.DialectWrapper, res *systemType.ApigwRoute) *goqu.InsertDataset {
		return d.Insert(apigwRouteTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"endpoint":   res.Endpoint,
				"method":     res.Method,
				"enabled":    res.Enabled,
				"meta":       res.Meta,
				"rel_group":  res.Group,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			})
	}

	// apigwRouteUpsertQuery assembles (insert+on-conflict) query for replacing apigwRoutes
	//
	// This function is auto-generated
	apigwRouteUpsertQuery = func(d goqu.DialectWrapper, res *systemType.ApigwRoute) *goqu.InsertDataset {
		var target = `,id`

		return apigwRouteInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"endpoint":   res.Endpoint,
						"method":     res.Method,
						"enabled":    res.Enabled,
						"meta":       res.Meta,
						"rel_group":  res.Group,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
						"created_by": res.CreatedBy,
						"updated_by": res.UpdatedBy,
						"deleted_by": res.DeletedBy,
					},
				),
			)
	}

	// apigwRouteUpdateQuery assembles query for updating apigwRoutes
	//
	// This function is auto-generated
	apigwRouteUpdateQuery = func(d goqu.DialectWrapper, res *systemType.ApigwRoute) *goqu.UpdateDataset {
		return d.Update(apigwRouteTable).
			Set(goqu.Record{
				"endpoint":   res.Endpoint,
				"method":     res.Method,
				"enabled":    res.Enabled,
				"meta":       res.Meta,
				"rel_group":  res.Group,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			}).
			Where(apigwRoutePrimaryKeys(res))
	}

	// apigwRouteDeleteQuery assembles delete query for removing apigwRoutes
	//
	// This function is auto-generated
	apigwRouteDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(apigwRouteTable).Where(ee...)
	}

	// apigwRouteDeleteQuery assembles delete query for removing apigwRoutes
	//
	// This function is auto-generated
	apigwRouteTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(apigwRouteTable)
	}

	// apigwRoutePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	apigwRoutePrimaryKeys = func(res *systemType.ApigwRoute) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// applicationTable represents applications store table
	//
	// This value is auto-generated
	applicationTable = goqu.T("applications")

	// applicationSelectQuery assembles select query for fetching applications
	//
	// This function is auto-generated
	applicationSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"name",
			"enabled",
			"weight",
			"unify",
			"rel_owner",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(applicationTable)
	}

	// applicationInsertQuery assembles query inserting applications
	//
	// This function is auto-generated
	applicationInsertQuery = func(d goqu.DialectWrapper, res *systemType.Application) *goqu.InsertDataset {
		return d.Insert(applicationTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"name":       res.Name,
				"enabled":    res.Enabled,
				"weight":     res.Weight,
				"unify":      res.Unify,
				"rel_owner":  res.OwnerID,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
			})
	}

	// applicationUpsertQuery assembles (insert+on-conflict) query for replacing applications
	//
	// This function is auto-generated
	applicationUpsertQuery = func(d goqu.DialectWrapper, res *systemType.Application) *goqu.InsertDataset {
		var target = `,id`

		return applicationInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"name":       res.Name,
						"enabled":    res.Enabled,
						"weight":     res.Weight,
						"unify":      res.Unify,
						"rel_owner":  res.OwnerID,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
					},
				),
			)
	}

	// applicationUpdateQuery assembles query for updating applications
	//
	// This function is auto-generated
	applicationUpdateQuery = func(d goqu.DialectWrapper, res *systemType.Application) *goqu.UpdateDataset {
		return d.Update(applicationTable).
			Set(goqu.Record{
				"name":       res.Name,
				"enabled":    res.Enabled,
				"weight":     res.Weight,
				"unify":      res.Unify,
				"rel_owner":  res.OwnerID,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
			}).
			Where(applicationPrimaryKeys(res))
	}

	// applicationDeleteQuery assembles delete query for removing applications
	//
	// This function is auto-generated
	applicationDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(applicationTable).Where(ee...)
	}

	// applicationDeleteQuery assembles delete query for removing applications
	//
	// This function is auto-generated
	applicationTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(applicationTable)
	}

	// applicationPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	applicationPrimaryKeys = func(res *systemType.Application) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// attachmentTable represents attachments store table
	//
	// This value is auto-generated
	attachmentTable = goqu.T("attachments")

	// attachmentSelectQuery assembles select query for fetching attachments
	//
	// This function is auto-generated
	attachmentSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_owner",
			"kind",
			"url",
			"preview_url",
			"name",
			"meta",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(attachmentTable)
	}

	// attachmentInsertQuery assembles query inserting attachments
	//
	// This function is auto-generated
	attachmentInsertQuery = func(d goqu.DialectWrapper, res *systemType.Attachment) *goqu.InsertDataset {
		return d.Insert(attachmentTable).
			Rows(goqu.Record{
				"id":          res.ID,
				"rel_owner":   res.OwnerID,
				"kind":        res.Kind,
				"url":         res.Url,
				"preview_url": res.PreviewUrl,
				"name":        res.Name,
				"meta":        res.Meta,
				"created_at":  res.CreatedAt,
				"updated_at":  res.UpdatedAt,
				"deleted_at":  res.DeletedAt,
			})
	}

	// attachmentUpsertQuery assembles (insert+on-conflict) query for replacing attachments
	//
	// This function is auto-generated
	attachmentUpsertQuery = func(d goqu.DialectWrapper, res *systemType.Attachment) *goqu.InsertDataset {
		var target = `,id`

		return attachmentInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_owner":   res.OwnerID,
						"kind":        res.Kind,
						"url":         res.Url,
						"preview_url": res.PreviewUrl,
						"name":        res.Name,
						"meta":        res.Meta,
						"created_at":  res.CreatedAt,
						"updated_at":  res.UpdatedAt,
						"deleted_at":  res.DeletedAt,
					},
				),
			)
	}

	// attachmentUpdateQuery assembles query for updating attachments
	//
	// This function is auto-generated
	attachmentUpdateQuery = func(d goqu.DialectWrapper, res *systemType.Attachment) *goqu.UpdateDataset {
		return d.Update(attachmentTable).
			Set(goqu.Record{
				"rel_owner":   res.OwnerID,
				"kind":        res.Kind,
				"url":         res.Url,
				"preview_url": res.PreviewUrl,
				"name":        res.Name,
				"meta":        res.Meta,
				"created_at":  res.CreatedAt,
				"updated_at":  res.UpdatedAt,
				"deleted_at":  res.DeletedAt,
			}).
			Where(attachmentPrimaryKeys(res))
	}

	// attachmentDeleteQuery assembles delete query for removing attachments
	//
	// This function is auto-generated
	attachmentDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(attachmentTable).Where(ee...)
	}

	// attachmentDeleteQuery assembles delete query for removing attachments
	//
	// This function is auto-generated
	attachmentTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(attachmentTable)
	}

	// attachmentPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	attachmentPrimaryKeys = func(res *systemType.Attachment) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// authClientTable represents authClients store table
	//
	// This value is auto-generated
	authClientTable = goqu.T("auth_clients")

	// authClientSelectQuery assembles select query for fetching authClients
	//
	// This function is auto-generated
	authClientSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"meta",
			"secret",
			"scope",
			"valid_grant",
			"redirect_uri",
			"enabled",
			"trusted",
			"valid_from",
			"expires_at",
			"security",
			"owned_by",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(authClientTable)
	}

	// authClientInsertQuery assembles query inserting authClients
	//
	// This function is auto-generated
	authClientInsertQuery = func(d goqu.DialectWrapper, res *systemType.AuthClient) *goqu.InsertDataset {
		return d.Insert(authClientTable).
			Rows(goqu.Record{
				"id":           res.ID,
				"handle":       res.Handle,
				"meta":         res.Meta,
				"secret":       res.Secret,
				"scope":        res.Scope,
				"valid_grant":  res.ValidGrant,
				"redirect_uri": res.RedirectURI,
				"enabled":      res.Enabled,
				"trusted":      res.Trusted,
				"valid_from":   res.ValidFrom,
				"expires_at":   res.ExpiresAt,
				"security":     res.Security,
				"owned_by":     res.OwnedBy,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
				"created_by":   res.CreatedBy,
				"updated_by":   res.UpdatedBy,
				"deleted_by":   res.DeletedBy,
			})
	}

	// authClientUpsertQuery assembles (insert+on-conflict) query for replacing authClients
	//
	// This function is auto-generated
	authClientUpsertQuery = func(d goqu.DialectWrapper, res *systemType.AuthClient) *goqu.InsertDataset {
		var target = `,id`

		return authClientInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":       res.Handle,
						"meta":         res.Meta,
						"secret":       res.Secret,
						"scope":        res.Scope,
						"valid_grant":  res.ValidGrant,
						"redirect_uri": res.RedirectURI,
						"enabled":      res.Enabled,
						"trusted":      res.Trusted,
						"valid_from":   res.ValidFrom,
						"expires_at":   res.ExpiresAt,
						"security":     res.Security,
						"owned_by":     res.OwnedBy,
						"created_at":   res.CreatedAt,
						"updated_at":   res.UpdatedAt,
						"deleted_at":   res.DeletedAt,
						"created_by":   res.CreatedBy,
						"updated_by":   res.UpdatedBy,
						"deleted_by":   res.DeletedBy,
					},
				),
			)
	}

	// authClientUpdateQuery assembles query for updating authClients
	//
	// This function is auto-generated
	authClientUpdateQuery = func(d goqu.DialectWrapper, res *systemType.AuthClient) *goqu.UpdateDataset {
		return d.Update(authClientTable).
			Set(goqu.Record{
				"handle":       res.Handle,
				"meta":         res.Meta,
				"secret":       res.Secret,
				"scope":        res.Scope,
				"valid_grant":  res.ValidGrant,
				"redirect_uri": res.RedirectURI,
				"enabled":      res.Enabled,
				"trusted":      res.Trusted,
				"valid_from":   res.ValidFrom,
				"expires_at":   res.ExpiresAt,
				"security":     res.Security,
				"owned_by":     res.OwnedBy,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
				"created_by":   res.CreatedBy,
				"updated_by":   res.UpdatedBy,
				"deleted_by":   res.DeletedBy,
			}).
			Where(authClientPrimaryKeys(res))
	}

	// authClientDeleteQuery assembles delete query for removing authClients
	//
	// This function is auto-generated
	authClientDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(authClientTable).Where(ee...)
	}

	// authClientDeleteQuery assembles delete query for removing authClients
	//
	// This function is auto-generated
	authClientTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(authClientTable)
	}

	// authClientPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	authClientPrimaryKeys = func(res *systemType.AuthClient) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// authConfirmedClientTable represents authConfirmedClients store table
	//
	// This value is auto-generated
	authConfirmedClientTable = goqu.T("auth_confirmed_clients")

	// authConfirmedClientSelectQuery assembles select query for fetching authConfirmedClients
	//
	// This function is auto-generated
	authConfirmedClientSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"rel_user",
			"rel_client",
			"confirmed_at",
		).From(authConfirmedClientTable)
	}

	// authConfirmedClientInsertQuery assembles query inserting authConfirmedClients
	//
	// This function is auto-generated
	authConfirmedClientInsertQuery = func(d goqu.DialectWrapper, res *systemType.AuthConfirmedClient) *goqu.InsertDataset {
		return d.Insert(authConfirmedClientTable).
			Rows(goqu.Record{
				"rel_user":     res.UserID,
				"rel_client":   res.ClientID,
				"confirmed_at": res.ConfirmedAt,
			})
	}

	// authConfirmedClientUpsertQuery assembles (insert+on-conflict) query for replacing authConfirmedClients
	//
	// This function is auto-generated
	authConfirmedClientUpsertQuery = func(d goqu.DialectWrapper, res *systemType.AuthConfirmedClient) *goqu.InsertDataset {
		var target = `,rel_user,rel_client`

		return authConfirmedClientInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"confirmed_at": res.ConfirmedAt,
					},
				),
			)
	}

	// authConfirmedClientUpdateQuery assembles query for updating authConfirmedClients
	//
	// This function is auto-generated
	authConfirmedClientUpdateQuery = func(d goqu.DialectWrapper, res *systemType.AuthConfirmedClient) *goqu.UpdateDataset {
		return d.Update(authConfirmedClientTable).
			Set(goqu.Record{
				"confirmed_at": res.ConfirmedAt,
			}).
			Where(authConfirmedClientPrimaryKeys(res))
	}

	// authConfirmedClientDeleteQuery assembles delete query for removing authConfirmedClients
	//
	// This function is auto-generated
	authConfirmedClientDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(authConfirmedClientTable).Where(ee...)
	}

	// authConfirmedClientDeleteQuery assembles delete query for removing authConfirmedClients
	//
	// This function is auto-generated
	authConfirmedClientTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(authConfirmedClientTable)
	}

	// authConfirmedClientPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	authConfirmedClientPrimaryKeys = func(res *systemType.AuthConfirmedClient) goqu.Ex {
		return goqu.Ex{
			"rel_user":   res.UserID,
			"rel_client": res.ClientID,
		}
	}

	// authOa2tokenTable represents authOa2tokens store table
	//
	// This value is auto-generated
	authOa2tokenTable = goqu.T("auth_oa2tokens")

	// authOa2tokenSelectQuery assembles select query for fetching authOa2tokens
	//
	// This function is auto-generated
	authOa2tokenSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"code",
			"access",
			"refresh",
			"data",
			"remote_addr",
			"user_agent",
			"rel_client",
			"rel_user",
			"expires_at",
			"created_at",
		).From(authOa2tokenTable)
	}

	// authOa2tokenInsertQuery assembles query inserting authOa2tokens
	//
	// This function is auto-generated
	authOa2tokenInsertQuery = func(d goqu.DialectWrapper, res *systemType.AuthOa2token) *goqu.InsertDataset {
		return d.Insert(authOa2tokenTable).
			Rows(goqu.Record{
				"id":          res.ID,
				"code":        res.Code,
				"access":      res.Access,
				"refresh":     res.Refresh,
				"data":        res.Data,
				"remote_addr": res.RemoteAddr,
				"user_agent":  res.UserAgent,
				"rel_client":  res.ClientID,
				"rel_user":    res.UserID,
				"expires_at":  res.ExpiresAt,
				"created_at":  res.CreatedAt,
			})
	}

	// authOa2tokenUpsertQuery assembles (insert+on-conflict) query for replacing authOa2tokens
	//
	// This function is auto-generated
	authOa2tokenUpsertQuery = func(d goqu.DialectWrapper, res *systemType.AuthOa2token) *goqu.InsertDataset {
		var target = `,id`

		return authOa2tokenInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"code":        res.Code,
						"access":      res.Access,
						"refresh":     res.Refresh,
						"data":        res.Data,
						"remote_addr": res.RemoteAddr,
						"user_agent":  res.UserAgent,
						"rel_client":  res.ClientID,
						"rel_user":    res.UserID,
						"expires_at":  res.ExpiresAt,
						"created_at":  res.CreatedAt,
					},
				),
			)
	}

	// authOa2tokenUpdateQuery assembles query for updating authOa2tokens
	//
	// This function is auto-generated
	authOa2tokenUpdateQuery = func(d goqu.DialectWrapper, res *systemType.AuthOa2token) *goqu.UpdateDataset {
		return d.Update(authOa2tokenTable).
			Set(goqu.Record{
				"code":        res.Code,
				"access":      res.Access,
				"refresh":     res.Refresh,
				"data":        res.Data,
				"remote_addr": res.RemoteAddr,
				"user_agent":  res.UserAgent,
				"rel_client":  res.ClientID,
				"rel_user":    res.UserID,
				"expires_at":  res.ExpiresAt,
				"created_at":  res.CreatedAt,
			}).
			Where(authOa2tokenPrimaryKeys(res))
	}

	// authOa2tokenDeleteQuery assembles delete query for removing authOa2tokens
	//
	// This function is auto-generated
	authOa2tokenDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(authOa2tokenTable).Where(ee...)
	}

	// authOa2tokenDeleteQuery assembles delete query for removing authOa2tokens
	//
	// This function is auto-generated
	authOa2tokenTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(authOa2tokenTable)
	}

	// authOa2tokenPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	authOa2tokenPrimaryKeys = func(res *systemType.AuthOa2token) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// authSessionTable represents authSessions store table
	//
	// This value is auto-generated
	authSessionTable = goqu.T("auth_sessions")

	// authSessionSelectQuery assembles select query for fetching authSessions
	//
	// This function is auto-generated
	authSessionSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"data",
			"rel_user",
			"remote_addr",
			"user_agent",
			"expires_at",
			"created_at",
		).From(authSessionTable)
	}

	// authSessionInsertQuery assembles query inserting authSessions
	//
	// This function is auto-generated
	authSessionInsertQuery = func(d goqu.DialectWrapper, res *systemType.AuthSession) *goqu.InsertDataset {
		return d.Insert(authSessionTable).
			Rows(goqu.Record{
				"id":          res.ID,
				"data":        res.Data,
				"rel_user":    res.UserID,
				"remote_addr": res.RemoteAddr,
				"user_agent":  res.UserAgent,
				"expires_at":  res.ExpiresAt,
				"created_at":  res.CreatedAt,
			})
	}

	// authSessionUpsertQuery assembles (insert+on-conflict) query for replacing authSessions
	//
	// This function is auto-generated
	authSessionUpsertQuery = func(d goqu.DialectWrapper, res *systemType.AuthSession) *goqu.InsertDataset {
		var target = `,id`

		return authSessionInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"data":        res.Data,
						"rel_user":    res.UserID,
						"remote_addr": res.RemoteAddr,
						"user_agent":  res.UserAgent,
						"expires_at":  res.ExpiresAt,
						"created_at":  res.CreatedAt,
					},
				),
			)
	}

	// authSessionUpdateQuery assembles query for updating authSessions
	//
	// This function is auto-generated
	authSessionUpdateQuery = func(d goqu.DialectWrapper, res *systemType.AuthSession) *goqu.UpdateDataset {
		return d.Update(authSessionTable).
			Set(goqu.Record{
				"data":        res.Data,
				"rel_user":    res.UserID,
				"remote_addr": res.RemoteAddr,
				"user_agent":  res.UserAgent,
				"expires_at":  res.ExpiresAt,
				"created_at":  res.CreatedAt,
			}).
			Where(authSessionPrimaryKeys(res))
	}

	// authSessionDeleteQuery assembles delete query for removing authSessions
	//
	// This function is auto-generated
	authSessionDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(authSessionTable).Where(ee...)
	}

	// authSessionDeleteQuery assembles delete query for removing authSessions
	//
	// This function is auto-generated
	authSessionTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(authSessionTable)
	}

	// authSessionPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	authSessionPrimaryKeys = func(res *systemType.AuthSession) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// automationSessionTable represents automationSessions store table
	//
	// This value is auto-generated
	automationSessionTable = goqu.T("automation_sessions")

	// automationSessionSelectQuery assembles select query for fetching automationSessions
	//
	// This function is auto-generated
	automationSessionSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_workflow",
			"event_type",
			"resource_type",
			"status",
			"input",
			"output",
			"stacktrace",
			"created_by",
			"created_at",
			"purge_at",
			"completed_at",
			"suspended_at",
			"error",
		).From(automationSessionTable)
	}

	// automationSessionInsertQuery assembles query inserting automationSessions
	//
	// This function is auto-generated
	automationSessionInsertQuery = func(d goqu.DialectWrapper, res *automationType.Session) *goqu.InsertDataset {
		return d.Insert(automationSessionTable).
			Rows(goqu.Record{
				"id":            res.ID,
				"rel_workflow":  res.WorkflowID,
				"event_type":    res.EventType,
				"resource_type": res.ResourceType,
				"status":        res.Status,
				"input":         res.Input,
				"output":        res.Output,
				"stacktrace":    res.Stacktrace,
				"created_by":    res.CreatedBy,
				"created_at":    res.CreatedAt,
				"purge_at":      res.PurgeAt,
				"completed_at":  res.CompletedAt,
				"suspended_at":  res.SuspendedAt,
				"error":         res.Error,
			})
	}

	// automationSessionUpsertQuery assembles (insert+on-conflict) query for replacing automationSessions
	//
	// This function is auto-generated
	automationSessionUpsertQuery = func(d goqu.DialectWrapper, res *automationType.Session) *goqu.InsertDataset {
		var target = `,id`

		return automationSessionInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_workflow":  res.WorkflowID,
						"event_type":    res.EventType,
						"resource_type": res.ResourceType,
						"status":        res.Status,
						"input":         res.Input,
						"output":        res.Output,
						"stacktrace":    res.Stacktrace,
						"created_by":    res.CreatedBy,
						"created_at":    res.CreatedAt,
						"purge_at":      res.PurgeAt,
						"completed_at":  res.CompletedAt,
						"suspended_at":  res.SuspendedAt,
						"error":         res.Error,
					},
				),
			)
	}

	// automationSessionUpdateQuery assembles query for updating automationSessions
	//
	// This function is auto-generated
	automationSessionUpdateQuery = func(d goqu.DialectWrapper, res *automationType.Session) *goqu.UpdateDataset {
		return d.Update(automationSessionTable).
			Set(goqu.Record{
				"rel_workflow":  res.WorkflowID,
				"event_type":    res.EventType,
				"resource_type": res.ResourceType,
				"status":        res.Status,
				"input":         res.Input,
				"output":        res.Output,
				"stacktrace":    res.Stacktrace,
				"created_by":    res.CreatedBy,
				"created_at":    res.CreatedAt,
				"purge_at":      res.PurgeAt,
				"completed_at":  res.CompletedAt,
				"suspended_at":  res.SuspendedAt,
				"error":         res.Error,
			}).
			Where(automationSessionPrimaryKeys(res))
	}

	// automationSessionDeleteQuery assembles delete query for removing automationSessions
	//
	// This function is auto-generated
	automationSessionDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(automationSessionTable).Where(ee...)
	}

	// automationSessionDeleteQuery assembles delete query for removing automationSessions
	//
	// This function is auto-generated
	automationSessionTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(automationSessionTable)
	}

	// automationSessionPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	automationSessionPrimaryKeys = func(res *automationType.Session) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// automationTriggerTable represents automationTriggers store table
	//
	// This value is auto-generated
	automationTriggerTable = goqu.T("automation_triggers")

	// automationTriggerSelectQuery assembles select query for fetching automationTriggers
	//
	// This function is auto-generated
	automationTriggerSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_workflow",
			"rel_step",
			"enabled",
			"resource_type",
			"event_type",
			"meta",
			"constraints",
			"input",
			"created_at",
			"updated_at",
			"deleted_at",
			"owned_by",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(automationTriggerTable)
	}

	// automationTriggerInsertQuery assembles query inserting automationTriggers
	//
	// This function is auto-generated
	automationTriggerInsertQuery = func(d goqu.DialectWrapper, res *automationType.Trigger) *goqu.InsertDataset {
		return d.Insert(automationTriggerTable).
			Rows(goqu.Record{
				"id":            res.ID,
				"rel_workflow":  res.WorkflowID,
				"rel_step":      res.StepID,
				"enabled":       res.Enabled,
				"resource_type": res.ResourceType,
				"event_type":    res.EventType,
				"meta":          res.Meta,
				"constraints":   res.Constraints,
				"input":         res.Input,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
				"owned_by":      res.OwnedBy,
				"created_by":    res.CreatedBy,
				"updated_by":    res.UpdatedBy,
				"deleted_by":    res.DeletedBy,
			})
	}

	// automationTriggerUpsertQuery assembles (insert+on-conflict) query for replacing automationTriggers
	//
	// This function is auto-generated
	automationTriggerUpsertQuery = func(d goqu.DialectWrapper, res *automationType.Trigger) *goqu.InsertDataset {
		var target = `,id`

		return automationTriggerInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_workflow":  res.WorkflowID,
						"rel_step":      res.StepID,
						"enabled":       res.Enabled,
						"resource_type": res.ResourceType,
						"event_type":    res.EventType,
						"meta":          res.Meta,
						"constraints":   res.Constraints,
						"input":         res.Input,
						"created_at":    res.CreatedAt,
						"updated_at":    res.UpdatedAt,
						"deleted_at":    res.DeletedAt,
						"owned_by":      res.OwnedBy,
						"created_by":    res.CreatedBy,
						"updated_by":    res.UpdatedBy,
						"deleted_by":    res.DeletedBy,
					},
				),
			)
	}

	// automationTriggerUpdateQuery assembles query for updating automationTriggers
	//
	// This function is auto-generated
	automationTriggerUpdateQuery = func(d goqu.DialectWrapper, res *automationType.Trigger) *goqu.UpdateDataset {
		return d.Update(automationTriggerTable).
			Set(goqu.Record{
				"rel_workflow":  res.WorkflowID,
				"rel_step":      res.StepID,
				"enabled":       res.Enabled,
				"resource_type": res.ResourceType,
				"event_type":    res.EventType,
				"meta":          res.Meta,
				"constraints":   res.Constraints,
				"input":         res.Input,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
				"owned_by":      res.OwnedBy,
				"created_by":    res.CreatedBy,
				"updated_by":    res.UpdatedBy,
				"deleted_by":    res.DeletedBy,
			}).
			Where(automationTriggerPrimaryKeys(res))
	}

	// automationTriggerDeleteQuery assembles delete query for removing automationTriggers
	//
	// This function is auto-generated
	automationTriggerDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(automationTriggerTable).Where(ee...)
	}

	// automationTriggerDeleteQuery assembles delete query for removing automationTriggers
	//
	// This function is auto-generated
	automationTriggerTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(automationTriggerTable)
	}

	// automationTriggerPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	automationTriggerPrimaryKeys = func(res *automationType.Trigger) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// automationWorkflowTable represents automationWorkflows store table
	//
	// This value is auto-generated
	automationWorkflowTable = goqu.T("automation_workflows")

	// automationWorkflowSelectQuery assembles select query for fetching automationWorkflows
	//
	// This function is auto-generated
	automationWorkflowSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"meta",
			"enabled",
			"trace",
			"keep_sessions",
			"scope",
			"steps",
			"paths",
			"issues",
			"run_as",
			"created_at",
			"updated_at",
			"deleted_at",
			"owned_by",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(automationWorkflowTable)
	}

	// automationWorkflowInsertQuery assembles query inserting automationWorkflows
	//
	// This function is auto-generated
	automationWorkflowInsertQuery = func(d goqu.DialectWrapper, res *automationType.Workflow) *goqu.InsertDataset {
		return d.Insert(automationWorkflowTable).
			Rows(goqu.Record{
				"id":            res.ID,
				"handle":        res.Handle,
				"meta":          res.Meta,
				"enabled":       res.Enabled,
				"trace":         res.Trace,
				"keep_sessions": res.KeepSessions,
				"scope":         res.Scope,
				"steps":         res.Steps,
				"paths":         res.Paths,
				"issues":        res.Issues,
				"run_as":        res.RunAs,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
				"owned_by":      res.OwnedBy,
				"created_by":    res.CreatedBy,
				"updated_by":    res.UpdatedBy,
				"deleted_by":    res.DeletedBy,
			})
	}

	// automationWorkflowUpsertQuery assembles (insert+on-conflict) query for replacing automationWorkflows
	//
	// This function is auto-generated
	automationWorkflowUpsertQuery = func(d goqu.DialectWrapper, res *automationType.Workflow) *goqu.InsertDataset {
		var target = `,id`

		return automationWorkflowInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":        res.Handle,
						"meta":          res.Meta,
						"enabled":       res.Enabled,
						"trace":         res.Trace,
						"keep_sessions": res.KeepSessions,
						"scope":         res.Scope,
						"steps":         res.Steps,
						"paths":         res.Paths,
						"issues":        res.Issues,
						"run_as":        res.RunAs,
						"created_at":    res.CreatedAt,
						"updated_at":    res.UpdatedAt,
						"deleted_at":    res.DeletedAt,
						"owned_by":      res.OwnedBy,
						"created_by":    res.CreatedBy,
						"updated_by":    res.UpdatedBy,
						"deleted_by":    res.DeletedBy,
					},
				),
			)
	}

	// automationWorkflowUpdateQuery assembles query for updating automationWorkflows
	//
	// This function is auto-generated
	automationWorkflowUpdateQuery = func(d goqu.DialectWrapper, res *automationType.Workflow) *goqu.UpdateDataset {
		return d.Update(automationWorkflowTable).
			Set(goqu.Record{
				"handle":        res.Handle,
				"meta":          res.Meta,
				"enabled":       res.Enabled,
				"trace":         res.Trace,
				"keep_sessions": res.KeepSessions,
				"scope":         res.Scope,
				"steps":         res.Steps,
				"paths":         res.Paths,
				"issues":        res.Issues,
				"run_as":        res.RunAs,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
				"owned_by":      res.OwnedBy,
				"created_by":    res.CreatedBy,
				"updated_by":    res.UpdatedBy,
				"deleted_by":    res.DeletedBy,
			}).
			Where(automationWorkflowPrimaryKeys(res))
	}

	// automationWorkflowDeleteQuery assembles delete query for removing automationWorkflows
	//
	// This function is auto-generated
	automationWorkflowDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(automationWorkflowTable).Where(ee...)
	}

	// automationWorkflowDeleteQuery assembles delete query for removing automationWorkflows
	//
	// This function is auto-generated
	automationWorkflowTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(automationWorkflowTable)
	}

	// automationWorkflowPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	automationWorkflowPrimaryKeys = func(res *automationType.Workflow) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// composeAttachmentTable represents composeAttachments store table
	//
	// This value is auto-generated
	composeAttachmentTable = goqu.T("compose_attachment")

	// composeAttachmentSelectQuery assembles select query for fetching composeAttachments
	//
	// This function is auto-generated
	composeAttachmentSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_owner",
			"rel_namespace",
			"kind",
			"url",
			"preview_url",
			"name",
			"meta",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(composeAttachmentTable)
	}

	// composeAttachmentInsertQuery assembles query inserting composeAttachments
	//
	// This function is auto-generated
	composeAttachmentInsertQuery = func(d goqu.DialectWrapper, res *composeType.Attachment) *goqu.InsertDataset {
		return d.Insert(composeAttachmentTable).
			Rows(goqu.Record{
				"id":            res.ID,
				"rel_owner":     res.OwnerID,
				"rel_namespace": res.NamespaceID,
				"kind":          res.Kind,
				"url":           res.Url,
				"preview_url":   res.PreviewUrl,
				"name":          res.Name,
				"meta":          res.Meta,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			})
	}

	// composeAttachmentUpsertQuery assembles (insert+on-conflict) query for replacing composeAttachments
	//
	// This function is auto-generated
	composeAttachmentUpsertQuery = func(d goqu.DialectWrapper, res *composeType.Attachment) *goqu.InsertDataset {
		var target = `,id`

		return composeAttachmentInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_owner":     res.OwnerID,
						"rel_namespace": res.NamespaceID,
						"kind":          res.Kind,
						"url":           res.Url,
						"preview_url":   res.PreviewUrl,
						"name":          res.Name,
						"meta":          res.Meta,
						"created_at":    res.CreatedAt,
						"updated_at":    res.UpdatedAt,
						"deleted_at":    res.DeletedAt,
					},
				),
			)
	}

	// composeAttachmentUpdateQuery assembles query for updating composeAttachments
	//
	// This function is auto-generated
	composeAttachmentUpdateQuery = func(d goqu.DialectWrapper, res *composeType.Attachment) *goqu.UpdateDataset {
		return d.Update(composeAttachmentTable).
			Set(goqu.Record{
				"rel_owner":     res.OwnerID,
				"rel_namespace": res.NamespaceID,
				"kind":          res.Kind,
				"url":           res.Url,
				"preview_url":   res.PreviewUrl,
				"name":          res.Name,
				"meta":          res.Meta,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			}).
			Where(composeAttachmentPrimaryKeys(res))
	}

	// composeAttachmentDeleteQuery assembles delete query for removing composeAttachments
	//
	// This function is auto-generated
	composeAttachmentDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(composeAttachmentTable).Where(ee...)
	}

	// composeAttachmentDeleteQuery assembles delete query for removing composeAttachments
	//
	// This function is auto-generated
	composeAttachmentTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(composeAttachmentTable)
	}

	// composeAttachmentPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	composeAttachmentPrimaryKeys = func(res *composeType.Attachment) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// composeChartTable represents composeCharts store table
	//
	// This value is auto-generated
	composeChartTable = goqu.T("compose_chart")

	// composeChartSelectQuery assembles select query for fetching composeCharts
	//
	// This function is auto-generated
	composeChartSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"name",
			"config",
			"rel_namespace",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(composeChartTable)
	}

	// composeChartInsertQuery assembles query inserting composeCharts
	//
	// This function is auto-generated
	composeChartInsertQuery = func(d goqu.DialectWrapper, res *composeType.Chart) *goqu.InsertDataset {
		return d.Insert(composeChartTable).
			Rows(goqu.Record{
				"id":            res.ID,
				"handle":        res.Handle,
				"name":          res.Name,
				"config":        res.Config,
				"rel_namespace": res.NamespaceID,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			})
	}

	// composeChartUpsertQuery assembles (insert+on-conflict) query for replacing composeCharts
	//
	// This function is auto-generated
	composeChartUpsertQuery = func(d goqu.DialectWrapper, res *composeType.Chart) *goqu.InsertDataset {
		var target = `,id`

		return composeChartInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":        res.Handle,
						"name":          res.Name,
						"config":        res.Config,
						"rel_namespace": res.NamespaceID,
						"created_at":    res.CreatedAt,
						"updated_at":    res.UpdatedAt,
						"deleted_at":    res.DeletedAt,
					},
				),
			)
	}

	// composeChartUpdateQuery assembles query for updating composeCharts
	//
	// This function is auto-generated
	composeChartUpdateQuery = func(d goqu.DialectWrapper, res *composeType.Chart) *goqu.UpdateDataset {
		return d.Update(composeChartTable).
			Set(goqu.Record{
				"handle":        res.Handle,
				"name":          res.Name,
				"config":        res.Config,
				"rel_namespace": res.NamespaceID,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			}).
			Where(composeChartPrimaryKeys(res))
	}

	// composeChartDeleteQuery assembles delete query for removing composeCharts
	//
	// This function is auto-generated
	composeChartDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(composeChartTable).Where(ee...)
	}

	// composeChartDeleteQuery assembles delete query for removing composeCharts
	//
	// This function is auto-generated
	composeChartTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(composeChartTable)
	}

	// composeChartPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	composeChartPrimaryKeys = func(res *composeType.Chart) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// composeModuleTable represents composeModules store table
	//
	// This value is auto-generated
	composeModuleTable = goqu.T("compose_module")

	// composeModuleSelectQuery assembles select query for fetching composeModules
	//
	// This function is auto-generated
	composeModuleSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"meta",
			"config",
			"rel_namespace",
			"name",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(composeModuleTable)
	}

	// composeModuleInsertQuery assembles query inserting composeModules
	//
	// This function is auto-generated
	composeModuleInsertQuery = func(d goqu.DialectWrapper, res *composeType.Module) *goqu.InsertDataset {
		return d.Insert(composeModuleTable).
			Rows(goqu.Record{
				"id":            res.ID,
				"handle":        res.Handle,
				"meta":          res.Meta,
				"config":        res.Config,
				"rel_namespace": res.NamespaceID,
				"name":          res.Name,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			})
	}

	// composeModuleUpsertQuery assembles (insert+on-conflict) query for replacing composeModules
	//
	// This function is auto-generated
	composeModuleUpsertQuery = func(d goqu.DialectWrapper, res *composeType.Module) *goqu.InsertDataset {
		var target = `,id`

		return composeModuleInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":        res.Handle,
						"meta":          res.Meta,
						"config":        res.Config,
						"rel_namespace": res.NamespaceID,
						"name":          res.Name,
						"created_at":    res.CreatedAt,
						"updated_at":    res.UpdatedAt,
						"deleted_at":    res.DeletedAt,
					},
				),
			)
	}

	// composeModuleUpdateQuery assembles query for updating composeModules
	//
	// This function is auto-generated
	composeModuleUpdateQuery = func(d goqu.DialectWrapper, res *composeType.Module) *goqu.UpdateDataset {
		return d.Update(composeModuleTable).
			Set(goqu.Record{
				"handle":        res.Handle,
				"meta":          res.Meta,
				"config":        res.Config,
				"rel_namespace": res.NamespaceID,
				"name":          res.Name,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			}).
			Where(composeModulePrimaryKeys(res))
	}

	// composeModuleDeleteQuery assembles delete query for removing composeModules
	//
	// This function is auto-generated
	composeModuleDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(composeModuleTable).Where(ee...)
	}

	// composeModuleDeleteQuery assembles delete query for removing composeModules
	//
	// This function is auto-generated
	composeModuleTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(composeModuleTable)
	}

	// composeModulePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	composeModulePrimaryKeys = func(res *composeType.Module) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// composeModuleFieldTable represents composeModuleFields store table
	//
	// This value is auto-generated
	composeModuleFieldTable = goqu.T("compose_module_field")

	// composeModuleFieldSelectQuery assembles select query for fetching composeModuleFields
	//
	// This function is auto-generated
	composeModuleFieldSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_module",
			"place",
			"kind",
			"name",
			"label",
			"options",
			"config",
			"is_required",
			"is_multi",
			"default_value",
			"expressions",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(composeModuleFieldTable)
	}

	// composeModuleFieldInsertQuery assembles query inserting composeModuleFields
	//
	// This function is auto-generated
	composeModuleFieldInsertQuery = func(d goqu.DialectWrapper, res *composeType.ModuleField) *goqu.InsertDataset {
		return d.Insert(composeModuleFieldTable).
			Rows(goqu.Record{
				"id":            res.ID,
				"rel_module":    res.ModuleID,
				"place":         res.Place,
				"kind":          res.Kind,
				"name":          res.Name,
				"label":         res.Label,
				"options":       res.Options,
				"config":        res.Config,
				"is_required":   res.Required,
				"is_multi":      res.Multi,
				"default_value": res.DefaultValue,
				"expressions":   res.Expressions,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			})
	}

	// composeModuleFieldUpsertQuery assembles (insert+on-conflict) query for replacing composeModuleFields
	//
	// This function is auto-generated
	composeModuleFieldUpsertQuery = func(d goqu.DialectWrapper, res *composeType.ModuleField) *goqu.InsertDataset {
		var target = `,id`

		return composeModuleFieldInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_module":    res.ModuleID,
						"place":         res.Place,
						"kind":          res.Kind,
						"name":          res.Name,
						"label":         res.Label,
						"options":       res.Options,
						"config":        res.Config,
						"is_required":   res.Required,
						"is_multi":      res.Multi,
						"default_value": res.DefaultValue,
						"expressions":   res.Expressions,
						"created_at":    res.CreatedAt,
						"updated_at":    res.UpdatedAt,
						"deleted_at":    res.DeletedAt,
					},
				),
			)
	}

	// composeModuleFieldUpdateQuery assembles query for updating composeModuleFields
	//
	// This function is auto-generated
	composeModuleFieldUpdateQuery = func(d goqu.DialectWrapper, res *composeType.ModuleField) *goqu.UpdateDataset {
		return d.Update(composeModuleFieldTable).
			Set(goqu.Record{
				"rel_module":    res.ModuleID,
				"place":         res.Place,
				"kind":          res.Kind,
				"name":          res.Name,
				"label":         res.Label,
				"options":       res.Options,
				"config":        res.Config,
				"is_required":   res.Required,
				"is_multi":      res.Multi,
				"default_value": res.DefaultValue,
				"expressions":   res.Expressions,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			}).
			Where(composeModuleFieldPrimaryKeys(res))
	}

	// composeModuleFieldDeleteQuery assembles delete query for removing composeModuleFields
	//
	// This function is auto-generated
	composeModuleFieldDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(composeModuleFieldTable).Where(ee...)
	}

	// composeModuleFieldDeleteQuery assembles delete query for removing composeModuleFields
	//
	// This function is auto-generated
	composeModuleFieldTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(composeModuleFieldTable)
	}

	// composeModuleFieldPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	composeModuleFieldPrimaryKeys = func(res *composeType.ModuleField) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// composeNamespaceTable represents composeNamespaces store table
	//
	// This value is auto-generated
	composeNamespaceTable = goqu.T("compose_namespace")

	// composeNamespaceSelectQuery assembles select query for fetching composeNamespaces
	//
	// This function is auto-generated
	composeNamespaceSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"slug",
			"enabled",
			"meta",
			"name",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(composeNamespaceTable)
	}

	// composeNamespaceInsertQuery assembles query inserting composeNamespaces
	//
	// This function is auto-generated
	composeNamespaceInsertQuery = func(d goqu.DialectWrapper, res *composeType.Namespace) *goqu.InsertDataset {
		return d.Insert(composeNamespaceTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"slug":       res.Slug,
				"enabled":    res.Enabled,
				"meta":       res.Meta,
				"name":       res.Name,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
			})
	}

	// composeNamespaceUpsertQuery assembles (insert+on-conflict) query for replacing composeNamespaces
	//
	// This function is auto-generated
	composeNamespaceUpsertQuery = func(d goqu.DialectWrapper, res *composeType.Namespace) *goqu.InsertDataset {
		var target = `,id`

		return composeNamespaceInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"slug":       res.Slug,
						"enabled":    res.Enabled,
						"meta":       res.Meta,
						"name":       res.Name,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
					},
				),
			)
	}

	// composeNamespaceUpdateQuery assembles query for updating composeNamespaces
	//
	// This function is auto-generated
	composeNamespaceUpdateQuery = func(d goqu.DialectWrapper, res *composeType.Namespace) *goqu.UpdateDataset {
		return d.Update(composeNamespaceTable).
			Set(goqu.Record{
				"slug":       res.Slug,
				"enabled":    res.Enabled,
				"meta":       res.Meta,
				"name":       res.Name,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
			}).
			Where(composeNamespacePrimaryKeys(res))
	}

	// composeNamespaceDeleteQuery assembles delete query for removing composeNamespaces
	//
	// This function is auto-generated
	composeNamespaceDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(composeNamespaceTable).Where(ee...)
	}

	// composeNamespaceDeleteQuery assembles delete query for removing composeNamespaces
	//
	// This function is auto-generated
	composeNamespaceTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(composeNamespaceTable)
	}

	// composeNamespacePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	composeNamespacePrimaryKeys = func(res *composeType.Namespace) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// composePageTable represents composePages store table
	//
	// This value is auto-generated
	composePageTable = goqu.T("compose_page")

	// composePageSelectQuery assembles select query for fetching composePages
	//
	// This function is auto-generated
	composePageSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"self_id",
			"rel_module",
			"rel_namespace",
			"handle",
			"config",
			"blocks",
			"visible",
			"weight",
			"title",
			"description",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(composePageTable)
	}

	// composePageInsertQuery assembles query inserting composePages
	//
	// This function is auto-generated
	composePageInsertQuery = func(d goqu.DialectWrapper, res *composeType.Page) *goqu.InsertDataset {
		return d.Insert(composePageTable).
			Rows(goqu.Record{
				"id":            res.ID,
				"self_id":       res.SelfID,
				"rel_module":    res.ModuleID,
				"rel_namespace": res.NamespaceID,
				"handle":        res.Handle,
				"config":        res.Config,
				"blocks":        res.Blocks,
				"visible":       res.Visible,
				"weight":        res.Weight,
				"title":         res.Title,
				"description":   res.Description,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			})
	}

	// composePageUpsertQuery assembles (insert+on-conflict) query for replacing composePages
	//
	// This function is auto-generated
	composePageUpsertQuery = func(d goqu.DialectWrapper, res *composeType.Page) *goqu.InsertDataset {
		var target = `,id`

		return composePageInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"self_id":       res.SelfID,
						"rel_module":    res.ModuleID,
						"rel_namespace": res.NamespaceID,
						"handle":        res.Handle,
						"config":        res.Config,
						"blocks":        res.Blocks,
						"visible":       res.Visible,
						"weight":        res.Weight,
						"title":         res.Title,
						"description":   res.Description,
						"created_at":    res.CreatedAt,
						"updated_at":    res.UpdatedAt,
						"deleted_at":    res.DeletedAt,
					},
				),
			)
	}

	// composePageUpdateQuery assembles query for updating composePages
	//
	// This function is auto-generated
	composePageUpdateQuery = func(d goqu.DialectWrapper, res *composeType.Page) *goqu.UpdateDataset {
		return d.Update(composePageTable).
			Set(goqu.Record{
				"self_id":       res.SelfID,
				"rel_module":    res.ModuleID,
				"rel_namespace": res.NamespaceID,
				"handle":        res.Handle,
				"config":        res.Config,
				"blocks":        res.Blocks,
				"visible":       res.Visible,
				"weight":        res.Weight,
				"title":         res.Title,
				"description":   res.Description,
				"created_at":    res.CreatedAt,
				"updated_at":    res.UpdatedAt,
				"deleted_at":    res.DeletedAt,
			}).
			Where(composePagePrimaryKeys(res))
	}

	// composePageDeleteQuery assembles delete query for removing composePages
	//
	// This function is auto-generated
	composePageDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(composePageTable).Where(ee...)
	}

	// composePageDeleteQuery assembles delete query for removing composePages
	//
	// This function is auto-generated
	composePageTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(composePageTable)
	}

	// composePagePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	composePagePrimaryKeys = func(res *composeType.Page) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// credentialTable represents credentials store table
	//
	// This value is auto-generated
	credentialTable = goqu.T("credentials")

	// credentialSelectQuery assembles select query for fetching credentials
	//
	// This function is auto-generated
	credentialSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_owner",
			"label",
			"kind",
			"credentials",
			"meta",
			"created_at",
			"updated_at",
			"deleted_at",
			"last_used_at",
			"expires_at",
		).From(credentialTable)
	}

	// credentialInsertQuery assembles query inserting credentials
	//
	// This function is auto-generated
	credentialInsertQuery = func(d goqu.DialectWrapper, res *systemType.Credential) *goqu.InsertDataset {
		return d.Insert(credentialTable).
			Rows(goqu.Record{
				"id":           res.ID,
				"rel_owner":    res.OwnerID,
				"label":        res.Label,
				"kind":         res.Kind,
				"credentials":  res.Credentials,
				"meta":         res.Meta,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
				"last_used_at": res.LastUsedAt,
				"expires_at":   res.ExpiresAt,
			})
	}

	// credentialUpsertQuery assembles (insert+on-conflict) query for replacing credentials
	//
	// This function is auto-generated
	credentialUpsertQuery = func(d goqu.DialectWrapper, res *systemType.Credential) *goqu.InsertDataset {
		var target = `,id`

		return credentialInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_owner":    res.OwnerID,
						"label":        res.Label,
						"kind":         res.Kind,
						"credentials":  res.Credentials,
						"meta":         res.Meta,
						"created_at":   res.CreatedAt,
						"updated_at":   res.UpdatedAt,
						"deleted_at":   res.DeletedAt,
						"last_used_at": res.LastUsedAt,
						"expires_at":   res.ExpiresAt,
					},
				),
			)
	}

	// credentialUpdateQuery assembles query for updating credentials
	//
	// This function is auto-generated
	credentialUpdateQuery = func(d goqu.DialectWrapper, res *systemType.Credential) *goqu.UpdateDataset {
		return d.Update(credentialTable).
			Set(goqu.Record{
				"rel_owner":    res.OwnerID,
				"label":        res.Label,
				"kind":         res.Kind,
				"credentials":  res.Credentials,
				"meta":         res.Meta,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
				"last_used_at": res.LastUsedAt,
				"expires_at":   res.ExpiresAt,
			}).
			Where(credentialPrimaryKeys(res))
	}

	// credentialDeleteQuery assembles delete query for removing credentials
	//
	// This function is auto-generated
	credentialDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(credentialTable).Where(ee...)
	}

	// credentialDeleteQuery assembles delete query for removing credentials
	//
	// This function is auto-generated
	credentialTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(credentialTable)
	}

	// credentialPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	credentialPrimaryKeys = func(res *systemType.Credential) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// dalConnectionTable represents dalConnections store table
	//
	// This value is auto-generated
	dalConnectionTable = goqu.T("dal_connections")

	// dalConnectionSelectQuery assembles select query for fetching dalConnections
	//
	// This function is auto-generated
	dalConnectionSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"type",
			"config",
			"meta",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(dalConnectionTable)
	}

	// dalConnectionInsertQuery assembles query inserting dalConnections
	//
	// This function is auto-generated
	dalConnectionInsertQuery = func(d goqu.DialectWrapper, res *systemType.DalConnection) *goqu.InsertDataset {
		return d.Insert(dalConnectionTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"handle":     res.Handle,
				"type":       res.Type,
				"config":     res.Config,
				"meta":       res.Meta,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			})
	}

	// dalConnectionUpsertQuery assembles (insert+on-conflict) query for replacing dalConnections
	//
	// This function is auto-generated
	dalConnectionUpsertQuery = func(d goqu.DialectWrapper, res *systemType.DalConnection) *goqu.InsertDataset {
		var target = `,id`

		return dalConnectionInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":     res.Handle,
						"type":       res.Type,
						"config":     res.Config,
						"meta":       res.Meta,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
						"created_by": res.CreatedBy,
						"updated_by": res.UpdatedBy,
						"deleted_by": res.DeletedBy,
					},
				),
			)
	}

	// dalConnectionUpdateQuery assembles query for updating dalConnections
	//
	// This function is auto-generated
	dalConnectionUpdateQuery = func(d goqu.DialectWrapper, res *systemType.DalConnection) *goqu.UpdateDataset {
		return d.Update(dalConnectionTable).
			Set(goqu.Record{
				"handle":     res.Handle,
				"type":       res.Type,
				"config":     res.Config,
				"meta":       res.Meta,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			}).
			Where(dalConnectionPrimaryKeys(res))
	}

	// dalConnectionDeleteQuery assembles delete query for removing dalConnections
	//
	// This function is auto-generated
	dalConnectionDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(dalConnectionTable).Where(ee...)
	}

	// dalConnectionDeleteQuery assembles delete query for removing dalConnections
	//
	// This function is auto-generated
	dalConnectionTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(dalConnectionTable)
	}

	// dalConnectionPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	dalConnectionPrimaryKeys = func(res *systemType.DalConnection) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// dalSensitivityLevelTable represents dalSensitivityLevels store table
	//
	// This value is auto-generated
	dalSensitivityLevelTable = goqu.T("dal_sensitivity_levels")

	// dalSensitivityLevelSelectQuery assembles select query for fetching dalSensitivityLevels
	//
	// This function is auto-generated
	dalSensitivityLevelSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"level",
			"meta",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(dalSensitivityLevelTable)
	}

	// dalSensitivityLevelInsertQuery assembles query inserting dalSensitivityLevels
	//
	// This function is auto-generated
	dalSensitivityLevelInsertQuery = func(d goqu.DialectWrapper, res *systemType.DalSensitivityLevel) *goqu.InsertDataset {
		return d.Insert(dalSensitivityLevelTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"handle":     res.Handle,
				"level":      res.Level,
				"meta":       res.Meta,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			})
	}

	// dalSensitivityLevelUpsertQuery assembles (insert+on-conflict) query for replacing dalSensitivityLevels
	//
	// This function is auto-generated
	dalSensitivityLevelUpsertQuery = func(d goqu.DialectWrapper, res *systemType.DalSensitivityLevel) *goqu.InsertDataset {
		var target = `,id`

		return dalSensitivityLevelInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":     res.Handle,
						"level":      res.Level,
						"meta":       res.Meta,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
						"created_by": res.CreatedBy,
						"updated_by": res.UpdatedBy,
						"deleted_by": res.DeletedBy,
					},
				),
			)
	}

	// dalSensitivityLevelUpdateQuery assembles query for updating dalSensitivityLevels
	//
	// This function is auto-generated
	dalSensitivityLevelUpdateQuery = func(d goqu.DialectWrapper, res *systemType.DalSensitivityLevel) *goqu.UpdateDataset {
		return d.Update(dalSensitivityLevelTable).
			Set(goqu.Record{
				"handle":     res.Handle,
				"level":      res.Level,
				"meta":       res.Meta,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			}).
			Where(dalSensitivityLevelPrimaryKeys(res))
	}

	// dalSensitivityLevelDeleteQuery assembles delete query for removing dalSensitivityLevels
	//
	// This function is auto-generated
	dalSensitivityLevelDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(dalSensitivityLevelTable).Where(ee...)
	}

	// dalSensitivityLevelDeleteQuery assembles delete query for removing dalSensitivityLevels
	//
	// This function is auto-generated
	dalSensitivityLevelTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(dalSensitivityLevelTable)
	}

	// dalSensitivityLevelPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	dalSensitivityLevelPrimaryKeys = func(res *systemType.DalSensitivityLevel) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// dataPrivacyRequestTable represents dataPrivacyRequests store table
	//
	// This value is auto-generated
	dataPrivacyRequestTable = goqu.T("data_privacy_requests")

	// dataPrivacyRequestSelectQuery assembles select query for fetching dataPrivacyRequests
	//
	// This function is auto-generated
	dataPrivacyRequestSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"kind",
			"status",
			"payload",
			"requested_at",
			"requested_by",
			"completed_at",
			"completed_by",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(dataPrivacyRequestTable)
	}

	// dataPrivacyRequestInsertQuery assembles query inserting dataPrivacyRequests
	//
	// This function is auto-generated
	dataPrivacyRequestInsertQuery = func(d goqu.DialectWrapper, res *systemType.DataPrivacyRequest) *goqu.InsertDataset {
		return d.Insert(dataPrivacyRequestTable).
			Rows(goqu.Record{
				"id":           res.ID,
				"kind":         res.Kind,
				"status":       res.Status,
				"payload":      res.Payload,
				"requested_at": res.RequestedAt,
				"requested_by": res.RequestedBy,
				"completed_at": res.CompletedAt,
				"completed_by": res.CompletedBy,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
				"created_by":   res.CreatedBy,
				"updated_by":   res.UpdatedBy,
				"deleted_by":   res.DeletedBy,
			})
	}

	// dataPrivacyRequestUpsertQuery assembles (insert+on-conflict) query for replacing dataPrivacyRequests
	//
	// This function is auto-generated
	dataPrivacyRequestUpsertQuery = func(d goqu.DialectWrapper, res *systemType.DataPrivacyRequest) *goqu.InsertDataset {
		var target = `,id`

		return dataPrivacyRequestInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"kind":         res.Kind,
						"status":       res.Status,
						"payload":      res.Payload,
						"requested_at": res.RequestedAt,
						"requested_by": res.RequestedBy,
						"completed_at": res.CompletedAt,
						"completed_by": res.CompletedBy,
						"created_at":   res.CreatedAt,
						"updated_at":   res.UpdatedAt,
						"deleted_at":   res.DeletedAt,
						"created_by":   res.CreatedBy,
						"updated_by":   res.UpdatedBy,
						"deleted_by":   res.DeletedBy,
					},
				),
			)
	}

	// dataPrivacyRequestUpdateQuery assembles query for updating dataPrivacyRequests
	//
	// This function is auto-generated
	dataPrivacyRequestUpdateQuery = func(d goqu.DialectWrapper, res *systemType.DataPrivacyRequest) *goqu.UpdateDataset {
		return d.Update(dataPrivacyRequestTable).
			Set(goqu.Record{
				"kind":         res.Kind,
				"status":       res.Status,
				"payload":      res.Payload,
				"requested_at": res.RequestedAt,
				"requested_by": res.RequestedBy,
				"completed_at": res.CompletedAt,
				"completed_by": res.CompletedBy,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
				"created_by":   res.CreatedBy,
				"updated_by":   res.UpdatedBy,
				"deleted_by":   res.DeletedBy,
			}).
			Where(dataPrivacyRequestPrimaryKeys(res))
	}

	// dataPrivacyRequestDeleteQuery assembles delete query for removing dataPrivacyRequests
	//
	// This function is auto-generated
	dataPrivacyRequestDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(dataPrivacyRequestTable).Where(ee...)
	}

	// dataPrivacyRequestDeleteQuery assembles delete query for removing dataPrivacyRequests
	//
	// This function is auto-generated
	dataPrivacyRequestTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(dataPrivacyRequestTable)
	}

	// dataPrivacyRequestPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	dataPrivacyRequestPrimaryKeys = func(res *systemType.DataPrivacyRequest) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// dataPrivacyRequestCommentTable represents dataPrivacyRequestComments store table
	//
	// This value is auto-generated
	dataPrivacyRequestCommentTable = goqu.T("data_privacy_request_comments")

	// dataPrivacyRequestCommentSelectQuery assembles select query for fetching dataPrivacyRequestComments
	//
	// This function is auto-generated
	dataPrivacyRequestCommentSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_request",
			"comment",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(dataPrivacyRequestCommentTable)
	}

	// dataPrivacyRequestCommentInsertQuery assembles query inserting dataPrivacyRequestComments
	//
	// This function is auto-generated
	dataPrivacyRequestCommentInsertQuery = func(d goqu.DialectWrapper, res *systemType.DataPrivacyRequestComment) *goqu.InsertDataset {
		return d.Insert(dataPrivacyRequestCommentTable).
			Rows(goqu.Record{
				"id":          res.ID,
				"rel_request": res.RequestID,
				"comment":     res.Comment,
				"created_at":  res.CreatedAt,
				"updated_at":  res.UpdatedAt,
				"deleted_at":  res.DeletedAt,
				"created_by":  res.CreatedBy,
				"updated_by":  res.UpdatedBy,
				"deleted_by":  res.DeletedBy,
			})
	}

	// dataPrivacyRequestCommentUpsertQuery assembles (insert+on-conflict) query for replacing dataPrivacyRequestComments
	//
	// This function is auto-generated
	dataPrivacyRequestCommentUpsertQuery = func(d goqu.DialectWrapper, res *systemType.DataPrivacyRequestComment) *goqu.InsertDataset {
		var target = `,id`

		return dataPrivacyRequestCommentInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_request": res.RequestID,
						"comment":     res.Comment,
						"created_at":  res.CreatedAt,
						"updated_at":  res.UpdatedAt,
						"deleted_at":  res.DeletedAt,
						"created_by":  res.CreatedBy,
						"updated_by":  res.UpdatedBy,
						"deleted_by":  res.DeletedBy,
					},
				),
			)
	}

	// dataPrivacyRequestCommentUpdateQuery assembles query for updating dataPrivacyRequestComments
	//
	// This function is auto-generated
	dataPrivacyRequestCommentUpdateQuery = func(d goqu.DialectWrapper, res *systemType.DataPrivacyRequestComment) *goqu.UpdateDataset {
		return d.Update(dataPrivacyRequestCommentTable).
			Set(goqu.Record{
				"rel_request": res.RequestID,
				"comment":     res.Comment,
				"created_at":  res.CreatedAt,
				"updated_at":  res.UpdatedAt,
				"deleted_at":  res.DeletedAt,
				"created_by":  res.CreatedBy,
				"updated_by":  res.UpdatedBy,
				"deleted_by":  res.DeletedBy,
			}).
			Where(dataPrivacyRequestCommentPrimaryKeys(res))
	}

	// dataPrivacyRequestCommentDeleteQuery assembles delete query for removing dataPrivacyRequestComments
	//
	// This function is auto-generated
	dataPrivacyRequestCommentDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(dataPrivacyRequestCommentTable).Where(ee...)
	}

	// dataPrivacyRequestCommentDeleteQuery assembles delete query for removing dataPrivacyRequestComments
	//
	// This function is auto-generated
	dataPrivacyRequestCommentTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(dataPrivacyRequestCommentTable)
	}

	// dataPrivacyRequestCommentPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	dataPrivacyRequestCommentPrimaryKeys = func(res *systemType.DataPrivacyRequestComment) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// federationExposedModuleTable represents federationExposedModules store table
	//
	// This value is auto-generated
	federationExposedModuleTable = goqu.T("federation_module_exposed")

	// federationExposedModuleSelectQuery assembles select query for fetching federationExposedModules
	//
	// This function is auto-generated
	federationExposedModuleSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"name",
			"rel_node",
			"rel_compose_module",
			"rel_compose_namespace",
			"fields",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(federationExposedModuleTable)
	}

	// federationExposedModuleInsertQuery assembles query inserting federationExposedModules
	//
	// This function is auto-generated
	federationExposedModuleInsertQuery = func(d goqu.DialectWrapper, res *federationType.ExposedModule) *goqu.InsertDataset {
		return d.Insert(federationExposedModuleTable).
			Rows(goqu.Record{
				"id":                    res.ID,
				"handle":                res.Handle,
				"name":                  res.Name,
				"rel_node":              res.NodeID,
				"rel_compose_module":    res.ComposeModuleID,
				"rel_compose_namespace": res.ComposeNamespaceID,
				"fields":                res.Fields,
				"created_at":            res.CreatedAt,
				"updated_at":            res.UpdatedAt,
				"deleted_at":            res.DeletedAt,
				"created_by":            res.CreatedBy,
				"updated_by":            res.UpdatedBy,
				"deleted_by":            res.DeletedBy,
			})
	}

	// federationExposedModuleUpsertQuery assembles (insert+on-conflict) query for replacing federationExposedModules
	//
	// This function is auto-generated
	federationExposedModuleUpsertQuery = func(d goqu.DialectWrapper, res *federationType.ExposedModule) *goqu.InsertDataset {
		var target = `,id`

		return federationExposedModuleInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":                res.Handle,
						"name":                  res.Name,
						"rel_node":              res.NodeID,
						"rel_compose_module":    res.ComposeModuleID,
						"rel_compose_namespace": res.ComposeNamespaceID,
						"fields":                res.Fields,
						"created_at":            res.CreatedAt,
						"updated_at":            res.UpdatedAt,
						"deleted_at":            res.DeletedAt,
						"created_by":            res.CreatedBy,
						"updated_by":            res.UpdatedBy,
						"deleted_by":            res.DeletedBy,
					},
				),
			)
	}

	// federationExposedModuleUpdateQuery assembles query for updating federationExposedModules
	//
	// This function is auto-generated
	federationExposedModuleUpdateQuery = func(d goqu.DialectWrapper, res *federationType.ExposedModule) *goqu.UpdateDataset {
		return d.Update(federationExposedModuleTable).
			Set(goqu.Record{
				"handle":                res.Handle,
				"name":                  res.Name,
				"rel_node":              res.NodeID,
				"rel_compose_module":    res.ComposeModuleID,
				"rel_compose_namespace": res.ComposeNamespaceID,
				"fields":                res.Fields,
				"created_at":            res.CreatedAt,
				"updated_at":            res.UpdatedAt,
				"deleted_at":            res.DeletedAt,
				"created_by":            res.CreatedBy,
				"updated_by":            res.UpdatedBy,
				"deleted_by":            res.DeletedBy,
			}).
			Where(federationExposedModulePrimaryKeys(res))
	}

	// federationExposedModuleDeleteQuery assembles delete query for removing federationExposedModules
	//
	// This function is auto-generated
	federationExposedModuleDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(federationExposedModuleTable).Where(ee...)
	}

	// federationExposedModuleDeleteQuery assembles delete query for removing federationExposedModules
	//
	// This function is auto-generated
	federationExposedModuleTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(federationExposedModuleTable)
	}

	// federationExposedModulePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	federationExposedModulePrimaryKeys = func(res *federationType.ExposedModule) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// federationModuleMappingTable represents federationModuleMappings store table
	//
	// This value is auto-generated
	federationModuleMappingTable = goqu.T("federation_module_mapping")

	// federationModuleMappingSelectQuery assembles select query for fetching federationModuleMappings
	//
	// This function is auto-generated
	federationModuleMappingSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"node_id",
			"federation_module_id",
			"compose_module_id",
			"compose_namespace_id",
			"field_mapping",
		).From(federationModuleMappingTable)
	}

	// federationModuleMappingInsertQuery assembles query inserting federationModuleMappings
	//
	// This function is auto-generated
	federationModuleMappingInsertQuery = func(d goqu.DialectWrapper, res *federationType.ModuleMapping) *goqu.InsertDataset {
		return d.Insert(federationModuleMappingTable).
			Rows(goqu.Record{
				"node_id":              res.NodeID,
				"federation_module_id": res.FederationModuleID,
				"compose_module_id":    res.ComposeModuleID,
				"compose_namespace_id": res.ComposeNamespaceID,
				"field_mapping":        res.FieldMapping,
			})
	}

	// federationModuleMappingUpsertQuery assembles (insert+on-conflict) query for replacing federationModuleMappings
	//
	// This function is auto-generated
	federationModuleMappingUpsertQuery = func(d goqu.DialectWrapper, res *federationType.ModuleMapping) *goqu.InsertDataset {
		var target = `,node_id`

		return federationModuleMappingInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"federation_module_id": res.FederationModuleID,
						"compose_module_id":    res.ComposeModuleID,
						"compose_namespace_id": res.ComposeNamespaceID,
						"field_mapping":        res.FieldMapping,
					},
				),
			)
	}

	// federationModuleMappingUpdateQuery assembles query for updating federationModuleMappings
	//
	// This function is auto-generated
	federationModuleMappingUpdateQuery = func(d goqu.DialectWrapper, res *federationType.ModuleMapping) *goqu.UpdateDataset {
		return d.Update(federationModuleMappingTable).
			Set(goqu.Record{
				"federation_module_id": res.FederationModuleID,
				"compose_module_id":    res.ComposeModuleID,
				"compose_namespace_id": res.ComposeNamespaceID,
				"field_mapping":        res.FieldMapping,
			}).
			Where(federationModuleMappingPrimaryKeys(res))
	}

	// federationModuleMappingDeleteQuery assembles delete query for removing federationModuleMappings
	//
	// This function is auto-generated
	federationModuleMappingDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(federationModuleMappingTable).Where(ee...)
	}

	// federationModuleMappingDeleteQuery assembles delete query for removing federationModuleMappings
	//
	// This function is auto-generated
	federationModuleMappingTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(federationModuleMappingTable)
	}

	// federationModuleMappingPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	federationModuleMappingPrimaryKeys = func(res *federationType.ModuleMapping) goqu.Ex {
		return goqu.Ex{
			"node_id": res.NodeID,
		}
	}

	// federationNodeTable represents federationNodes store table
	//
	// This value is auto-generated
	federationNodeTable = goqu.T("federation_nodes")

	// federationNodeSelectQuery assembles select query for fetching federationNodes
	//
	// This function is auto-generated
	federationNodeSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"name",
			"shared_node_id",
			"base_url",
			"status",
			"contact",
			"pair_token",
			"auth_token",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(federationNodeTable)
	}

	// federationNodeInsertQuery assembles query inserting federationNodes
	//
	// This function is auto-generated
	federationNodeInsertQuery = func(d goqu.DialectWrapper, res *federationType.Node) *goqu.InsertDataset {
		return d.Insert(federationNodeTable).
			Rows(goqu.Record{
				"id":             res.ID,
				"name":           res.Name,
				"shared_node_id": res.SharedNodeID,
				"base_url":       res.BaseURL,
				"status":         res.Status,
				"contact":        res.Contact,
				"pair_token":     res.PairToken,
				"auth_token":     res.AuthToken,
				"created_at":     res.CreatedAt,
				"updated_at":     res.UpdatedAt,
				"deleted_at":     res.DeletedAt,
				"created_by":     res.CreatedBy,
				"updated_by":     res.UpdatedBy,
				"deleted_by":     res.DeletedBy,
			})
	}

	// federationNodeUpsertQuery assembles (insert+on-conflict) query for replacing federationNodes
	//
	// This function is auto-generated
	federationNodeUpsertQuery = func(d goqu.DialectWrapper, res *federationType.Node) *goqu.InsertDataset {
		var target = `,id`

		return federationNodeInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"name":           res.Name,
						"shared_node_id": res.SharedNodeID,
						"base_url":       res.BaseURL,
						"status":         res.Status,
						"contact":        res.Contact,
						"pair_token":     res.PairToken,
						"auth_token":     res.AuthToken,
						"created_at":     res.CreatedAt,
						"updated_at":     res.UpdatedAt,
						"deleted_at":     res.DeletedAt,
						"created_by":     res.CreatedBy,
						"updated_by":     res.UpdatedBy,
						"deleted_by":     res.DeletedBy,
					},
				),
			)
	}

	// federationNodeUpdateQuery assembles query for updating federationNodes
	//
	// This function is auto-generated
	federationNodeUpdateQuery = func(d goqu.DialectWrapper, res *federationType.Node) *goqu.UpdateDataset {
		return d.Update(federationNodeTable).
			Set(goqu.Record{
				"name":           res.Name,
				"shared_node_id": res.SharedNodeID,
				"base_url":       res.BaseURL,
				"status":         res.Status,
				"contact":        res.Contact,
				"pair_token":     res.PairToken,
				"auth_token":     res.AuthToken,
				"created_at":     res.CreatedAt,
				"updated_at":     res.UpdatedAt,
				"deleted_at":     res.DeletedAt,
				"created_by":     res.CreatedBy,
				"updated_by":     res.UpdatedBy,
				"deleted_by":     res.DeletedBy,
			}).
			Where(federationNodePrimaryKeys(res))
	}

	// federationNodeDeleteQuery assembles delete query for removing federationNodes
	//
	// This function is auto-generated
	federationNodeDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(federationNodeTable).Where(ee...)
	}

	// federationNodeDeleteQuery assembles delete query for removing federationNodes
	//
	// This function is auto-generated
	federationNodeTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(federationNodeTable)
	}

	// federationNodePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	federationNodePrimaryKeys = func(res *federationType.Node) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// federationNodeSyncTable represents federationNodeSyncs store table
	//
	// This value is auto-generated
	federationNodeSyncTable = goqu.T("federation_nodes_sync")

	// federationNodeSyncSelectQuery assembles select query for fetching federationNodeSyncs
	//
	// This function is auto-generated
	federationNodeSyncSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"node_id",
			"module_id",
			"sync_type",
			"sync_status",
			"time_of_action",
		).From(federationNodeSyncTable)
	}

	// federationNodeSyncInsertQuery assembles query inserting federationNodeSyncs
	//
	// This function is auto-generated
	federationNodeSyncInsertQuery = func(d goqu.DialectWrapper, res *federationType.NodeSync) *goqu.InsertDataset {
		return d.Insert(federationNodeSyncTable).
			Rows(goqu.Record{
				"node_id":        res.NodeID,
				"module_id":      res.ModuleID,
				"sync_type":      res.SyncType,
				"sync_status":    res.SyncStatus,
				"time_of_action": res.TimeOfAction,
			})
	}

	// federationNodeSyncUpsertQuery assembles (insert+on-conflict) query for replacing federationNodeSyncs
	//
	// This function is auto-generated
	federationNodeSyncUpsertQuery = func(d goqu.DialectWrapper, res *federationType.NodeSync) *goqu.InsertDataset {
		var target = `,node_id`

		return federationNodeSyncInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"module_id":      res.ModuleID,
						"sync_type":      res.SyncType,
						"sync_status":    res.SyncStatus,
						"time_of_action": res.TimeOfAction,
					},
				),
			)
	}

	// federationNodeSyncUpdateQuery assembles query for updating federationNodeSyncs
	//
	// This function is auto-generated
	federationNodeSyncUpdateQuery = func(d goqu.DialectWrapper, res *federationType.NodeSync) *goqu.UpdateDataset {
		return d.Update(federationNodeSyncTable).
			Set(goqu.Record{
				"module_id":      res.ModuleID,
				"sync_type":      res.SyncType,
				"sync_status":    res.SyncStatus,
				"time_of_action": res.TimeOfAction,
			}).
			Where(federationNodeSyncPrimaryKeys(res))
	}

	// federationNodeSyncDeleteQuery assembles delete query for removing federationNodeSyncs
	//
	// This function is auto-generated
	federationNodeSyncDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(federationNodeSyncTable).Where(ee...)
	}

	// federationNodeSyncDeleteQuery assembles delete query for removing federationNodeSyncs
	//
	// This function is auto-generated
	federationNodeSyncTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(federationNodeSyncTable)
	}

	// federationNodeSyncPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	federationNodeSyncPrimaryKeys = func(res *federationType.NodeSync) goqu.Ex {
		return goqu.Ex{
			"node_id": res.NodeID,
		}
	}

	// federationSharedModuleTable represents federationSharedModules store table
	//
	// This value is auto-generated
	federationSharedModuleTable = goqu.T("federation_module_shared")

	// federationSharedModuleSelectQuery assembles select query for fetching federationSharedModules
	//
	// This function is auto-generated
	federationSharedModuleSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"rel_node",
			"name",
			"xref_module",
			"fields",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(federationSharedModuleTable)
	}

	// federationSharedModuleInsertQuery assembles query inserting federationSharedModules
	//
	// This function is auto-generated
	federationSharedModuleInsertQuery = func(d goqu.DialectWrapper, res *federationType.SharedModule) *goqu.InsertDataset {
		return d.Insert(federationSharedModuleTable).
			Rows(goqu.Record{
				"id":          res.ID,
				"handle":      res.Handle,
				"rel_node":    res.NodeID,
				"name":        res.Name,
				"xref_module": res.ExternalFederationModuleID,
				"fields":      res.Fields,
				"created_at":  res.CreatedAt,
				"updated_at":  res.UpdatedAt,
				"deleted_at":  res.DeletedAt,
				"created_by":  res.CreatedBy,
				"updated_by":  res.UpdatedBy,
				"deleted_by":  res.DeletedBy,
			})
	}

	// federationSharedModuleUpsertQuery assembles (insert+on-conflict) query for replacing federationSharedModules
	//
	// This function is auto-generated
	federationSharedModuleUpsertQuery = func(d goqu.DialectWrapper, res *federationType.SharedModule) *goqu.InsertDataset {
		var target = `,id`

		return federationSharedModuleInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":      res.Handle,
						"rel_node":    res.NodeID,
						"name":        res.Name,
						"xref_module": res.ExternalFederationModuleID,
						"fields":      res.Fields,
						"created_at":  res.CreatedAt,
						"updated_at":  res.UpdatedAt,
						"deleted_at":  res.DeletedAt,
						"created_by":  res.CreatedBy,
						"updated_by":  res.UpdatedBy,
						"deleted_by":  res.DeletedBy,
					},
				),
			)
	}

	// federationSharedModuleUpdateQuery assembles query for updating federationSharedModules
	//
	// This function is auto-generated
	federationSharedModuleUpdateQuery = func(d goqu.DialectWrapper, res *federationType.SharedModule) *goqu.UpdateDataset {
		return d.Update(federationSharedModuleTable).
			Set(goqu.Record{
				"handle":      res.Handle,
				"rel_node":    res.NodeID,
				"name":        res.Name,
				"xref_module": res.ExternalFederationModuleID,
				"fields":      res.Fields,
				"created_at":  res.CreatedAt,
				"updated_at":  res.UpdatedAt,
				"deleted_at":  res.DeletedAt,
				"created_by":  res.CreatedBy,
				"updated_by":  res.UpdatedBy,
				"deleted_by":  res.DeletedBy,
			}).
			Where(federationSharedModulePrimaryKeys(res))
	}

	// federationSharedModuleDeleteQuery assembles delete query for removing federationSharedModules
	//
	// This function is auto-generated
	federationSharedModuleDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(federationSharedModuleTable).Where(ee...)
	}

	// federationSharedModuleDeleteQuery assembles delete query for removing federationSharedModules
	//
	// This function is auto-generated
	federationSharedModuleTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(federationSharedModuleTable)
	}

	// federationSharedModulePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	federationSharedModulePrimaryKeys = func(res *federationType.SharedModule) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// flagTable represents flags store table
	//
	// This value is auto-generated
	flagTable = goqu.T("flags")

	// flagSelectQuery assembles select query for fetching flags
	//
	// This function is auto-generated
	flagSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"kind",
			"rel_resource",
			"owned_by",
			"name",
			"active",
		).From(flagTable)
	}

	// flagInsertQuery assembles query inserting flags
	//
	// This function is auto-generated
	flagInsertQuery = func(d goqu.DialectWrapper, res *flagType.Flag) *goqu.InsertDataset {
		return d.Insert(flagTable).
			Rows(goqu.Record{
				"kind":         res.Kind,
				"rel_resource": res.ResourceID,
				"owned_by":     res.OwnedBy,
				"name":         res.Name,
				"active":       res.Active,
			})
	}

	// flagUpsertQuery assembles (insert+on-conflict) query for replacing flags
	//
	// This function is auto-generated
	flagUpsertQuery = func(d goqu.DialectWrapper, res *flagType.Flag) *goqu.InsertDataset {
		var target = `,kind,rel_resource,owned_by,LOWER(name)`

		return flagInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"active": res.Active,
					},
				),
			)
	}

	// flagUpdateQuery assembles query for updating flags
	//
	// This function is auto-generated
	flagUpdateQuery = func(d goqu.DialectWrapper, res *flagType.Flag) *goqu.UpdateDataset {
		return d.Update(flagTable).
			Set(goqu.Record{
				"active": res.Active,
			}).
			Where(flagPrimaryKeys(res))
	}

	// flagDeleteQuery assembles delete query for removing flags
	//
	// This function is auto-generated
	flagDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(flagTable).Where(ee...)
	}

	// flagDeleteQuery assembles delete query for removing flags
	//
	// This function is auto-generated
	flagTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(flagTable)
	}

	// flagPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	flagPrimaryKeys = func(res *flagType.Flag) goqu.Ex {
		return goqu.Ex{
			"kind":         res.Kind,
			"rel_resource": res.ResourceID,
			"owned_by":     res.OwnedBy,
			"name":         res.Name,
		}
	}

	// labelTable represents labels store table
	//
	// This value is auto-generated
	labelTable = goqu.T("labels")

	// labelSelectQuery assembles select query for fetching labels
	//
	// This function is auto-generated
	labelSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"kind",
			"rel_resource",
			"name",
			"value",
		).From(labelTable)
	}

	// labelInsertQuery assembles query inserting labels
	//
	// This function is auto-generated
	labelInsertQuery = func(d goqu.DialectWrapper, res *labelsType.Label) *goqu.InsertDataset {
		return d.Insert(labelTable).
			Rows(goqu.Record{
				"kind":         res.Kind,
				"rel_resource": res.ResourceID,
				"name":         res.Name,
				"value":        res.Value,
			})
	}

	// labelUpsertQuery assembles (insert+on-conflict) query for replacing labels
	//
	// This function is auto-generated
	labelUpsertQuery = func(d goqu.DialectWrapper, res *labelsType.Label) *goqu.InsertDataset {
		var target = `,kind,rel_resource,LOWER(name)`

		return labelInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"value": res.Value,
					},
				),
			)
	}

	// labelUpdateQuery assembles query for updating labels
	//
	// This function is auto-generated
	labelUpdateQuery = func(d goqu.DialectWrapper, res *labelsType.Label) *goqu.UpdateDataset {
		return d.Update(labelTable).
			Set(goqu.Record{
				"value": res.Value,
			}).
			Where(labelPrimaryKeys(res))
	}

	// labelDeleteQuery assembles delete query for removing labels
	//
	// This function is auto-generated
	labelDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(labelTable).Where(ee...)
	}

	// labelDeleteQuery assembles delete query for removing labels
	//
	// This function is auto-generated
	labelTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(labelTable)
	}

	// labelPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	labelPrimaryKeys = func(res *labelsType.Label) goqu.Ex {
		return goqu.Ex{
			"kind":         res.Kind,
			"rel_resource": res.ResourceID,
			"name":         res.Name,
		}
	}

	// queueTable represents queues store table
	//
	// This value is auto-generated
	queueTable = goqu.T("queue_settings")

	// queueSelectQuery assembles select query for fetching queues
	//
	// This function is auto-generated
	queueSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"consumer",
			"queue",
			"meta",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(queueTable)
	}

	// queueInsertQuery assembles query inserting queues
	//
	// This function is auto-generated
	queueInsertQuery = func(d goqu.DialectWrapper, res *systemType.Queue) *goqu.InsertDataset {
		return d.Insert(queueTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"consumer":   res.Consumer,
				"queue":      res.Queue,
				"meta":       res.Meta,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			})
	}

	// queueUpsertQuery assembles (insert+on-conflict) query for replacing queues
	//
	// This function is auto-generated
	queueUpsertQuery = func(d goqu.DialectWrapper, res *systemType.Queue) *goqu.InsertDataset {
		var target = `,id`

		return queueInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"consumer":   res.Consumer,
						"queue":      res.Queue,
						"meta":       res.Meta,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
						"created_by": res.CreatedBy,
						"updated_by": res.UpdatedBy,
						"deleted_by": res.DeletedBy,
					},
				),
			)
	}

	// queueUpdateQuery assembles query for updating queues
	//
	// This function is auto-generated
	queueUpdateQuery = func(d goqu.DialectWrapper, res *systemType.Queue) *goqu.UpdateDataset {
		return d.Update(queueTable).
			Set(goqu.Record{
				"consumer":   res.Consumer,
				"queue":      res.Queue,
				"meta":       res.Meta,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			}).
			Where(queuePrimaryKeys(res))
	}

	// queueDeleteQuery assembles delete query for removing queues
	//
	// This function is auto-generated
	queueDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(queueTable).Where(ee...)
	}

	// queueDeleteQuery assembles delete query for removing queues
	//
	// This function is auto-generated
	queueTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(queueTable)
	}

	// queuePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	queuePrimaryKeys = func(res *systemType.Queue) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// queueMessageTable represents queueMessages store table
	//
	// This value is auto-generated
	queueMessageTable = goqu.T("queue_messages")

	// queueMessageSelectQuery assembles select query for fetching queueMessages
	//
	// This function is auto-generated
	queueMessageSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"queue",
			"payload",
			"processed",
			"created",
		).From(queueMessageTable)
	}

	// queueMessageInsertQuery assembles query inserting queueMessages
	//
	// This function is auto-generated
	queueMessageInsertQuery = func(d goqu.DialectWrapper, res *systemType.QueueMessage) *goqu.InsertDataset {
		return d.Insert(queueMessageTable).
			Rows(goqu.Record{
				"id":        res.ID,
				"queue":     res.Queue,
				"payload":   res.Payload,
				"processed": res.Processed,
				"created":   res.Created,
			})
	}

	// queueMessageUpsertQuery assembles (insert+on-conflict) query for replacing queueMessages
	//
	// This function is auto-generated
	queueMessageUpsertQuery = func(d goqu.DialectWrapper, res *systemType.QueueMessage) *goqu.InsertDataset {
		var target = `,id`

		return queueMessageInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"queue":     res.Queue,
						"payload":   res.Payload,
						"processed": res.Processed,
						"created":   res.Created,
					},
				),
			)
	}

	// queueMessageUpdateQuery assembles query for updating queueMessages
	//
	// This function is auto-generated
	queueMessageUpdateQuery = func(d goqu.DialectWrapper, res *systemType.QueueMessage) *goqu.UpdateDataset {
		return d.Update(queueMessageTable).
			Set(goqu.Record{
				"queue":     res.Queue,
				"payload":   res.Payload,
				"processed": res.Processed,
				"created":   res.Created,
			}).
			Where(queueMessagePrimaryKeys(res))
	}

	// queueMessageDeleteQuery assembles delete query for removing queueMessages
	//
	// This function is auto-generated
	queueMessageDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(queueMessageTable).Where(ee...)
	}

	// queueMessageDeleteQuery assembles delete query for removing queueMessages
	//
	// This function is auto-generated
	queueMessageTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(queueMessageTable)
	}

	// queueMessagePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	queueMessagePrimaryKeys = func(res *systemType.QueueMessage) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// rbacRuleTable represents rbacRules store table
	//
	// This value is auto-generated
	rbacRuleTable = goqu.T("rbac_rules")

	// rbacRuleSelectQuery assembles select query for fetching rbacRules
	//
	// This function is auto-generated
	rbacRuleSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"rel_role",
			"resource",
			"operation",
			"access",
		).From(rbacRuleTable)
	}

	// rbacRuleInsertQuery assembles query inserting rbacRules
	//
	// This function is auto-generated
	rbacRuleInsertQuery = func(d goqu.DialectWrapper, res *rbacType.Rule) *goqu.InsertDataset {
		return d.Insert(rbacRuleTable).
			Rows(goqu.Record{
				"rel_role":  res.RoleID,
				"resource":  res.Resource,
				"operation": res.Operation,
				"access":    res.Access,
			})
	}

	// rbacRuleUpsertQuery assembles (insert+on-conflict) query for replacing rbacRules
	//
	// This function is auto-generated
	rbacRuleUpsertQuery = func(d goqu.DialectWrapper, res *rbacType.Rule) *goqu.InsertDataset {
		var target = `,rel_role,resource,operation`

		return rbacRuleInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"access": res.Access,
					},
				),
			)
	}

	// rbacRuleUpdateQuery assembles query for updating rbacRules
	//
	// This function is auto-generated
	rbacRuleUpdateQuery = func(d goqu.DialectWrapper, res *rbacType.Rule) *goqu.UpdateDataset {
		return d.Update(rbacRuleTable).
			Set(goqu.Record{
				"access": res.Access,
			}).
			Where(rbacRulePrimaryKeys(res))
	}

	// rbacRuleDeleteQuery assembles delete query for removing rbacRules
	//
	// This function is auto-generated
	rbacRuleDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(rbacRuleTable).Where(ee...)
	}

	// rbacRuleDeleteQuery assembles delete query for removing rbacRules
	//
	// This function is auto-generated
	rbacRuleTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(rbacRuleTable)
	}

	// rbacRulePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	rbacRulePrimaryKeys = func(res *rbacType.Rule) goqu.Ex {
		return goqu.Ex{
			"rel_role":  res.RoleID,
			"resource":  res.Resource,
			"operation": res.Operation,
		}
	}

	// reminderTable represents reminders store table
	//
	// This value is auto-generated
	reminderTable = goqu.T("reminders")

	// reminderSelectQuery assembles select query for fetching reminders
	//
	// This function is auto-generated
	reminderSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"resource",
			"payload",
			"snooze_count",
			"assigned_to",
			"assigned_by",
			"assigned_at",
			"dismissed_by",
			"dismissed_at",
			"remind_at",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(reminderTable)
	}

	// reminderInsertQuery assembles query inserting reminders
	//
	// This function is auto-generated
	reminderInsertQuery = func(d goqu.DialectWrapper, res *systemType.Reminder) *goqu.InsertDataset {
		return d.Insert(reminderTable).
			Rows(goqu.Record{
				"id":           res.ID,
				"resource":     res.Resource,
				"payload":      res.Payload,
				"snooze_count": res.SnoozeCount,
				"assigned_to":  res.AssignedTo,
				"assigned_by":  res.AssignedBy,
				"assigned_at":  res.AssignedAt,
				"dismissed_by": res.DismissedBy,
				"dismissed_at": res.DismissedAt,
				"remind_at":    res.RemindAt,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
			})
	}

	// reminderUpsertQuery assembles (insert+on-conflict) query for replacing reminders
	//
	// This function is auto-generated
	reminderUpsertQuery = func(d goqu.DialectWrapper, res *systemType.Reminder) *goqu.InsertDataset {
		var target = `,id`

		return reminderInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"resource":     res.Resource,
						"payload":      res.Payload,
						"snooze_count": res.SnoozeCount,
						"assigned_to":  res.AssignedTo,
						"assigned_by":  res.AssignedBy,
						"assigned_at":  res.AssignedAt,
						"dismissed_by": res.DismissedBy,
						"dismissed_at": res.DismissedAt,
						"remind_at":    res.RemindAt,
						"created_at":   res.CreatedAt,
						"updated_at":   res.UpdatedAt,
						"deleted_at":   res.DeletedAt,
					},
				),
			)
	}

	// reminderUpdateQuery assembles query for updating reminders
	//
	// This function is auto-generated
	reminderUpdateQuery = func(d goqu.DialectWrapper, res *systemType.Reminder) *goqu.UpdateDataset {
		return d.Update(reminderTable).
			Set(goqu.Record{
				"resource":     res.Resource,
				"payload":      res.Payload,
				"snooze_count": res.SnoozeCount,
				"assigned_to":  res.AssignedTo,
				"assigned_by":  res.AssignedBy,
				"assigned_at":  res.AssignedAt,
				"dismissed_by": res.DismissedBy,
				"dismissed_at": res.DismissedAt,
				"remind_at":    res.RemindAt,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
			}).
			Where(reminderPrimaryKeys(res))
	}

	// reminderDeleteQuery assembles delete query for removing reminders
	//
	// This function is auto-generated
	reminderDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(reminderTable).Where(ee...)
	}

	// reminderDeleteQuery assembles delete query for removing reminders
	//
	// This function is auto-generated
	reminderTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(reminderTable)
	}

	// reminderPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	reminderPrimaryKeys = func(res *systemType.Reminder) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// reportTable represents reports store table
	//
	// This value is auto-generated
	reportTable = goqu.T("reports")

	// reportSelectQuery assembles select query for fetching reports
	//
	// This function is auto-generated
	reportSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"handle",
			"meta",
			"scenarios",
			"sources",
			"blocks",
			"owned_by",
			"created_at",
			"updated_at",
			"deleted_at",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(reportTable)
	}

	// reportInsertQuery assembles query inserting reports
	//
	// This function is auto-generated
	reportInsertQuery = func(d goqu.DialectWrapper, res *systemType.Report) *goqu.InsertDataset {
		return d.Insert(reportTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"handle":     res.Handle,
				"meta":       res.Meta,
				"scenarios":  res.Scenarios,
				"sources":    res.Sources,
				"blocks":     res.Blocks,
				"owned_by":   res.OwnedBy,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			})
	}

	// reportUpsertQuery assembles (insert+on-conflict) query for replacing reports
	//
	// This function is auto-generated
	reportUpsertQuery = func(d goqu.DialectWrapper, res *systemType.Report) *goqu.InsertDataset {
		var target = `,id`

		return reportInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"handle":     res.Handle,
						"meta":       res.Meta,
						"scenarios":  res.Scenarios,
						"sources":    res.Sources,
						"blocks":     res.Blocks,
						"owned_by":   res.OwnedBy,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
						"created_by": res.CreatedBy,
						"updated_by": res.UpdatedBy,
						"deleted_by": res.DeletedBy,
					},
				),
			)
	}

	// reportUpdateQuery assembles query for updating reports
	//
	// This function is auto-generated
	reportUpdateQuery = func(d goqu.DialectWrapper, res *systemType.Report) *goqu.UpdateDataset {
		return d.Update(reportTable).
			Set(goqu.Record{
				"handle":     res.Handle,
				"meta":       res.Meta,
				"scenarios":  res.Scenarios,
				"sources":    res.Sources,
				"blocks":     res.Blocks,
				"owned_by":   res.OwnedBy,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			}).
			Where(reportPrimaryKeys(res))
	}

	// reportDeleteQuery assembles delete query for removing reports
	//
	// This function is auto-generated
	reportDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(reportTable).Where(ee...)
	}

	// reportDeleteQuery assembles delete query for removing reports
	//
	// This function is auto-generated
	reportTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(reportTable)
	}

	// reportPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	reportPrimaryKeys = func(res *systemType.Report) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// resourceActivityTable represents resourceActivitys store table
	//
	// This value is auto-generated
	resourceActivityTable = goqu.T("resource_activity_log")

	// resourceActivitySelectQuery assembles select query for fetching resourceActivitys
	//
	// This function is auto-generated
	resourceActivitySelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"ts",
			"resource_type",
			"resource_action",
			"rel_resource",
			"meta",
		).From(resourceActivityTable)
	}

	// resourceActivityInsertQuery assembles query inserting resourceActivitys
	//
	// This function is auto-generated
	resourceActivityInsertQuery = func(d goqu.DialectWrapper, res *discoveryType.ResourceActivity) *goqu.InsertDataset {
		return d.Insert(resourceActivityTable).
			Rows(goqu.Record{
				"id":              res.ID,
				"ts":              res.Timestamp,
				"resource_type":   res.ResourceType,
				"resource_action": res.ResourceAction,
				"rel_resource":    res.ResourceID,
				"meta":            res.Meta,
			})
	}

	// resourceActivityUpsertQuery assembles (insert+on-conflict) query for replacing resourceActivitys
	//
	// This function is auto-generated
	resourceActivityUpsertQuery = func(d goqu.DialectWrapper, res *discoveryType.ResourceActivity) *goqu.InsertDataset {
		var target = `,id`

		return resourceActivityInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"ts":              res.Timestamp,
						"resource_type":   res.ResourceType,
						"resource_action": res.ResourceAction,
						"rel_resource":    res.ResourceID,
						"meta":            res.Meta,
					},
				),
			)
	}

	// resourceActivityUpdateQuery assembles query for updating resourceActivitys
	//
	// This function is auto-generated
	resourceActivityUpdateQuery = func(d goqu.DialectWrapper, res *discoveryType.ResourceActivity) *goqu.UpdateDataset {
		return d.Update(resourceActivityTable).
			Set(goqu.Record{
				"ts":              res.Timestamp,
				"resource_type":   res.ResourceType,
				"resource_action": res.ResourceAction,
				"rel_resource":    res.ResourceID,
				"meta":            res.Meta,
			}).
			Where(resourceActivityPrimaryKeys(res))
	}

	// resourceActivityDeleteQuery assembles delete query for removing resourceActivitys
	//
	// This function is auto-generated
	resourceActivityDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(resourceActivityTable).Where(ee...)
	}

	// resourceActivityDeleteQuery assembles delete query for removing resourceActivitys
	//
	// This function is auto-generated
	resourceActivityTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(resourceActivityTable)
	}

	// resourceActivityPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	resourceActivityPrimaryKeys = func(res *discoveryType.ResourceActivity) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// resourceTranslationTable represents resourceTranslations store table
	//
	// This value is auto-generated
	resourceTranslationTable = goqu.T("resource_translations")

	// resourceTranslationSelectQuery assembles select query for fetching resourceTranslations
	//
	// This function is auto-generated
	resourceTranslationSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"lang",
			"resource",
			"k",
			"message",
			"created_at",
			"updated_at",
			"deleted_at",
			"owned_by",
			"created_by",
			"updated_by",
			"deleted_by",
		).From(resourceTranslationTable)
	}

	// resourceTranslationInsertQuery assembles query inserting resourceTranslations
	//
	// This function is auto-generated
	resourceTranslationInsertQuery = func(d goqu.DialectWrapper, res *systemType.ResourceTranslation) *goqu.InsertDataset {
		return d.Insert(resourceTranslationTable).
			Rows(goqu.Record{
				"id":         res.ID,
				"lang":       res.Lang,
				"resource":   res.Resource,
				"k":          res.K,
				"message":    res.Message,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"owned_by":   res.OwnedBy,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			})
	}

	// resourceTranslationUpsertQuery assembles (insert+on-conflict) query for replacing resourceTranslations
	//
	// This function is auto-generated
	resourceTranslationUpsertQuery = func(d goqu.DialectWrapper, res *systemType.ResourceTranslation) *goqu.InsertDataset {
		var target = `,id`

		return resourceTranslationInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"lang":       res.Lang,
						"resource":   res.Resource,
						"k":          res.K,
						"message":    res.Message,
						"created_at": res.CreatedAt,
						"updated_at": res.UpdatedAt,
						"deleted_at": res.DeletedAt,
						"owned_by":   res.OwnedBy,
						"created_by": res.CreatedBy,
						"updated_by": res.UpdatedBy,
						"deleted_by": res.DeletedBy,
					},
				),
			)
	}

	// resourceTranslationUpdateQuery assembles query for updating resourceTranslations
	//
	// This function is auto-generated
	resourceTranslationUpdateQuery = func(d goqu.DialectWrapper, res *systemType.ResourceTranslation) *goqu.UpdateDataset {
		return d.Update(resourceTranslationTable).
			Set(goqu.Record{
				"lang":       res.Lang,
				"resource":   res.Resource,
				"k":          res.K,
				"message":    res.Message,
				"created_at": res.CreatedAt,
				"updated_at": res.UpdatedAt,
				"deleted_at": res.DeletedAt,
				"owned_by":   res.OwnedBy,
				"created_by": res.CreatedBy,
				"updated_by": res.UpdatedBy,
				"deleted_by": res.DeletedBy,
			}).
			Where(resourceTranslationPrimaryKeys(res))
	}

	// resourceTranslationDeleteQuery assembles delete query for removing resourceTranslations
	//
	// This function is auto-generated
	resourceTranslationDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(resourceTranslationTable).Where(ee...)
	}

	// resourceTranslationDeleteQuery assembles delete query for removing resourceTranslations
	//
	// This function is auto-generated
	resourceTranslationTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(resourceTranslationTable)
	}

	// resourceTranslationPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	resourceTranslationPrimaryKeys = func(res *systemType.ResourceTranslation) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// roleTable represents roles store table
	//
	// This value is auto-generated
	roleTable = goqu.T("roles")

	// roleSelectQuery assembles select query for fetching roles
	//
	// This function is auto-generated
	roleSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"name",
			"handle",
			"meta",
			"archived_at",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(roleTable)
	}

	// roleInsertQuery assembles query inserting roles
	//
	// This function is auto-generated
	roleInsertQuery = func(d goqu.DialectWrapper, res *systemType.Role) *goqu.InsertDataset {
		return d.Insert(roleTable).
			Rows(goqu.Record{
				"id":          res.ID,
				"name":        res.Name,
				"handle":      res.Handle,
				"meta":        res.Meta,
				"archived_at": res.ArchivedAt,
				"created_at":  res.CreatedAt,
				"updated_at":  res.UpdatedAt,
				"deleted_at":  res.DeletedAt,
			})
	}

	// roleUpsertQuery assembles (insert+on-conflict) query for replacing roles
	//
	// This function is auto-generated
	roleUpsertQuery = func(d goqu.DialectWrapper, res *systemType.Role) *goqu.InsertDataset {
		var target = `,id`

		return roleInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"name":        res.Name,
						"handle":      res.Handle,
						"meta":        res.Meta,
						"archived_at": res.ArchivedAt,
						"created_at":  res.CreatedAt,
						"updated_at":  res.UpdatedAt,
						"deleted_at":  res.DeletedAt,
					},
				),
			)
	}

	// roleUpdateQuery assembles query for updating roles
	//
	// This function is auto-generated
	roleUpdateQuery = func(d goqu.DialectWrapper, res *systemType.Role) *goqu.UpdateDataset {
		return d.Update(roleTable).
			Set(goqu.Record{
				"name":        res.Name,
				"handle":      res.Handle,
				"meta":        res.Meta,
				"archived_at": res.ArchivedAt,
				"created_at":  res.CreatedAt,
				"updated_at":  res.UpdatedAt,
				"deleted_at":  res.DeletedAt,
			}).
			Where(rolePrimaryKeys(res))
	}

	// roleDeleteQuery assembles delete query for removing roles
	//
	// This function is auto-generated
	roleDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(roleTable).Where(ee...)
	}

	// roleDeleteQuery assembles delete query for removing roles
	//
	// This function is auto-generated
	roleTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(roleTable)
	}

	// rolePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	rolePrimaryKeys = func(res *systemType.Role) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// roleMemberTable represents roleMembers store table
	//
	// This value is auto-generated
	roleMemberTable = goqu.T("role_members")

	// roleMemberSelectQuery assembles select query for fetching roleMembers
	//
	// This function is auto-generated
	roleMemberSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"rel_user",
			"rel_role",
		).From(roleMemberTable)
	}

	// roleMemberInsertQuery assembles query inserting roleMembers
	//
	// This function is auto-generated
	roleMemberInsertQuery = func(d goqu.DialectWrapper, res *systemType.RoleMember) *goqu.InsertDataset {
		return d.Insert(roleMemberTable).
			Rows(goqu.Record{
				"rel_user": res.UserID,
				"rel_role": res.RoleID,
			})
	}

	// roleMemberUpsertQuery assembles (insert+on-conflict) query for replacing roleMembers
	//
	// This function is auto-generated
	roleMemberUpsertQuery = func(d goqu.DialectWrapper, res *systemType.RoleMember) *goqu.InsertDataset {
		var target = `,rel_user,rel_role`

		return roleMemberInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{},
				),
			)
	}

	// roleMemberUpdateQuery assembles query for updating roleMembers
	//
	// This function is auto-generated
	roleMemberUpdateQuery = func(d goqu.DialectWrapper, res *systemType.RoleMember) *goqu.UpdateDataset {
		return d.Update(roleMemberTable).
			Set(goqu.Record{}).
			Where(roleMemberPrimaryKeys(res))
	}

	// roleMemberDeleteQuery assembles delete query for removing roleMembers
	//
	// This function is auto-generated
	roleMemberDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(roleMemberTable).Where(ee...)
	}

	// roleMemberDeleteQuery assembles delete query for removing roleMembers
	//
	// This function is auto-generated
	roleMemberTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(roleMemberTable)
	}

	// roleMemberPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	roleMemberPrimaryKeys = func(res *systemType.RoleMember) goqu.Ex {
		return goqu.Ex{
			"rel_user": res.UserID,
			"rel_role": res.RoleID,
		}
	}

	// settingValueTable represents settingValues store table
	//
	// This value is auto-generated
	settingValueTable = goqu.T("settings")

	// settingValueSelectQuery assembles select query for fetching settingValues
	//
	// This function is auto-generated
	settingValueSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"rel_owner",
			"name",
			"value",
			"updated_by",
			"updated_at",
		).From(settingValueTable)
	}

	// settingValueInsertQuery assembles query inserting settingValues
	//
	// This function is auto-generated
	settingValueInsertQuery = func(d goqu.DialectWrapper, res *systemType.SettingValue) *goqu.InsertDataset {
		return d.Insert(settingValueTable).
			Rows(goqu.Record{
				"rel_owner":  res.OwnedBy,
				"name":       res.Name,
				"value":      res.Value,
				"updated_by": res.UpdatedBy,
				"updated_at": res.UpdatedAt,
			})
	}

	// settingValueUpsertQuery assembles (insert+on-conflict) query for replacing settingValues
	//
	// This function is auto-generated
	settingValueUpsertQuery = func(d goqu.DialectWrapper, res *systemType.SettingValue) *goqu.InsertDataset {
		var target = `,rel_owner,LOWER(name)`

		return settingValueInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"value":      res.Value,
						"updated_by": res.UpdatedBy,
						"updated_at": res.UpdatedAt,
					},
				),
			)
	}

	// settingValueUpdateQuery assembles query for updating settingValues
	//
	// This function is auto-generated
	settingValueUpdateQuery = func(d goqu.DialectWrapper, res *systemType.SettingValue) *goqu.UpdateDataset {
		return d.Update(settingValueTable).
			Set(goqu.Record{
				"value":      res.Value,
				"updated_by": res.UpdatedBy,
				"updated_at": res.UpdatedAt,
			}).
			Where(settingValuePrimaryKeys(res))
	}

	// settingValueDeleteQuery assembles delete query for removing settingValues
	//
	// This function is auto-generated
	settingValueDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(settingValueTable).Where(ee...)
	}

	// settingValueDeleteQuery assembles delete query for removing settingValues
	//
	// This function is auto-generated
	settingValueTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(settingValueTable)
	}

	// settingValuePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	settingValuePrimaryKeys = func(res *systemType.SettingValue) goqu.Ex {
		return goqu.Ex{
			"rel_owner": res.OwnedBy,
			"name":      res.Name,
		}
	}

	// templateTable represents templates store table
	//
	// This value is auto-generated
	templateTable = goqu.T("templates")

	// templateSelectQuery assembles select query for fetching templates
	//
	// This function is auto-generated
	templateSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"rel_owner",
			"handle",
			"language",
			"type",
			"partial",
			"meta",
			"template",
			"created_at",
			"updated_at",
			"deleted_at",
			"last_used_at",
		).From(templateTable)
	}

	// templateInsertQuery assembles query inserting templates
	//
	// This function is auto-generated
	templateInsertQuery = func(d goqu.DialectWrapper, res *systemType.Template) *goqu.InsertDataset {
		return d.Insert(templateTable).
			Rows(goqu.Record{
				"id":           res.ID,
				"rel_owner":    res.OwnerID,
				"handle":       res.Handle,
				"language":     res.Language,
				"type":         res.Type,
				"partial":      res.Partial,
				"meta":         res.Meta,
				"template":     res.Template,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
				"last_used_at": res.LastUsedAt,
			})
	}

	// templateUpsertQuery assembles (insert+on-conflict) query for replacing templates
	//
	// This function is auto-generated
	templateUpsertQuery = func(d goqu.DialectWrapper, res *systemType.Template) *goqu.InsertDataset {
		var target = `,id`

		return templateInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"rel_owner":    res.OwnerID,
						"handle":       res.Handle,
						"language":     res.Language,
						"type":         res.Type,
						"partial":      res.Partial,
						"meta":         res.Meta,
						"template":     res.Template,
						"created_at":   res.CreatedAt,
						"updated_at":   res.UpdatedAt,
						"deleted_at":   res.DeletedAt,
						"last_used_at": res.LastUsedAt,
					},
				),
			)
	}

	// templateUpdateQuery assembles query for updating templates
	//
	// This function is auto-generated
	templateUpdateQuery = func(d goqu.DialectWrapper, res *systemType.Template) *goqu.UpdateDataset {
		return d.Update(templateTable).
			Set(goqu.Record{
				"rel_owner":    res.OwnerID,
				"handle":       res.Handle,
				"language":     res.Language,
				"type":         res.Type,
				"partial":      res.Partial,
				"meta":         res.Meta,
				"template":     res.Template,
				"created_at":   res.CreatedAt,
				"updated_at":   res.UpdatedAt,
				"deleted_at":   res.DeletedAt,
				"last_used_at": res.LastUsedAt,
			}).
			Where(templatePrimaryKeys(res))
	}

	// templateDeleteQuery assembles delete query for removing templates
	//
	// This function is auto-generated
	templateDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(templateTable).Where(ee...)
	}

	// templateDeleteQuery assembles delete query for removing templates
	//
	// This function is auto-generated
	templateTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(templateTable)
	}

	// templatePrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	templatePrimaryKeys = func(res *systemType.Template) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}

	// userTable represents users store table
	//
	// This value is auto-generated
	userTable = goqu.T("users")

	// userSelectQuery assembles select query for fetching users
	//
	// This function is auto-generated
	userSelectQuery = func(d goqu.DialectWrapper) *goqu.SelectDataset {
		return d.Select(
			"id",
			"email",
			"email_confirmed",
			"username",
			"name",
			"handle",
			"kind",
			"meta",
			"suspended_at",
			"created_at",
			"updated_at",
			"deleted_at",
		).From(userTable)
	}

	// userInsertQuery assembles query inserting users
	//
	// This function is auto-generated
	userInsertQuery = func(d goqu.DialectWrapper, res *systemType.User) *goqu.InsertDataset {
		return d.Insert(userTable).
			Rows(goqu.Record{
				"id":              res.ID,
				"email":           res.Email,
				"email_confirmed": res.EmailConfirmed,
				"username":        res.Username,
				"name":            res.Name,
				"handle":          res.Handle,
				"kind":            res.Kind,
				"meta":            res.Meta,
				"suspended_at":    res.SuspendedAt,
				"created_at":      res.CreatedAt,
				"updated_at":      res.UpdatedAt,
				"deleted_at":      res.DeletedAt,
			})
	}

	// userUpsertQuery assembles (insert+on-conflict) query for replacing users
	//
	// This function is auto-generated
	userUpsertQuery = func(d goqu.DialectWrapper, res *systemType.User) *goqu.InsertDataset {
		var target = `,id`

		return userInsertQuery(d, res).
			OnConflict(
				goqu.DoUpdate(target[1:],
					goqu.Record{
						"email":           res.Email,
						"email_confirmed": res.EmailConfirmed,
						"username":        res.Username,
						"name":            res.Name,
						"handle":          res.Handle,
						"kind":            res.Kind,
						"meta":            res.Meta,
						"suspended_at":    res.SuspendedAt,
						"created_at":      res.CreatedAt,
						"updated_at":      res.UpdatedAt,
						"deleted_at":      res.DeletedAt,
					},
				),
			)
	}

	// userUpdateQuery assembles query for updating users
	//
	// This function is auto-generated
	userUpdateQuery = func(d goqu.DialectWrapper, res *systemType.User) *goqu.UpdateDataset {
		return d.Update(userTable).
			Set(goqu.Record{
				"email":           res.Email,
				"email_confirmed": res.EmailConfirmed,
				"username":        res.Username,
				"name":            res.Name,
				"handle":          res.Handle,
				"kind":            res.Kind,
				"meta":            res.Meta,
				"suspended_at":    res.SuspendedAt,
				"created_at":      res.CreatedAt,
				"updated_at":      res.UpdatedAt,
				"deleted_at":      res.DeletedAt,
			}).
			Where(userPrimaryKeys(res))
	}

	// userDeleteQuery assembles delete query for removing users
	//
	// This function is auto-generated
	userDeleteQuery = func(d goqu.DialectWrapper, ee ...goqu.Expression) *goqu.DeleteDataset {
		return d.Delete(userTable).Where(ee...)
	}

	// userDeleteQuery assembles delete query for removing users
	//
	// This function is auto-generated
	userTruncateQuery = func(d goqu.DialectWrapper) *goqu.TruncateDataset {
		return d.Truncate(userTable)
	}

	// userPrimaryKeys assembles set of conditions for all primary keys
	//
	// This function is auto-generated
	userPrimaryKeys = func(res *systemType.User) goqu.Ex {
		return goqu.Ex{
			"id": res.ID,
		}
	}
)
