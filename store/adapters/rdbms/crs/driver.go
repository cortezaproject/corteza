package crs

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/qlng"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type (
	sqlizer interface {
		ToSQL() (string, []interface{}, error)
	}

	connection interface {
		sqlx.QueryerContext
		sqlx.ExecerContext
	}

	attrExpression interface {
		exp.Comparable
		exp.Inable
		exp.Isable
	}

	attributeType interface {
		Type() data.AttributeType
	}

	column struct {
		ident      string
		columnType attributeType
		attributes []*data.Attribute

		encode func([]*data.Attribute, *types.Record) (any, error)
		decode func([]*data.Attribute, any, *types.Record) error
	}

	queryParser interface {
		Parse(string) (exp.Expression, error)
	}

	crs struct {
		model *data.Model
		conn  connection

		queryParser queryParser

		deepJsonFn func(ident string, pp ...any) exp.LiteralExpression

		table exp.IdentifierExpression

		// ID column identifier expression
		sysColumnID exp.IdentifierExpression

		// optional record fields/columns/expressions
		sysExprNamespaceID attrExpression
		sysExprModuleID    attrExpression
		sysExprDeletedAt   attrExpression

		// all columns we're selecting from when
		// we're selecting from all columns
		columns []*column
	}
)

// CRS returns fully initialized model store
//
// It abstracts database table and its columns and provides unified interface
// for fetching and storing records.
func CRS(m *data.Model, c connection) *crs {
	var (
		ms = &crs{
			model: m,
			conn:  c,

			queryParser: ql.Converter(),
		}
	)

	ms.table = exp.NewIdentifierExpression("", m.Ident, "")
	ms.deepJsonFn = rdbms.DeepIdentJSON

	_ = ms.genColumns()

	ms.queryParser = ql.Converter(
		ql.SymHandler(func(node *qlng.ASTNode) (exp.Expression, error) {
			return ms.attrToExpr(node.Symbol)
		}),
	)

	return ms
}

// generates columns from model attributes.
//
// Important note!
//   Attribute name is not always equal to name of the column
//   (case of aliased fields or embedded values)
func (ms *crs) genColumns() error {
	var (
		colIdent string
		att      *data.Attribute

		uniqCols = make(map[string]int)
		cols     = make([]*column, 0, len(ms.model.Attributes))
	)

	for a := range ms.model.Attributes {
		att = ms.model.Attributes[a]
		colIdent = attrColumnIdent(att)

		// there are a couple of important system attributes/columns
		// we need to locate and link for easier management
		switch att.Ident {
		case sysID:
			ms.sysColumnID = ms.table.Col(colIdent)
		case sysNamespaceID:
			ms.sysExprNamespaceID = ms.table.Col(colIdent)
		case sysModuleID:
			ms.sysExprModuleID = ms.table.Col(colIdent)
		case sysDeletedAt:
			ms.sysExprDeletedAt = ms.table.Col(colIdent)
		}

		// identify db columns from list of attributes
		colIndex, has := uniqCols[colIdent]

		if has {
			// column already initialized one of the previous iterations
			// append attribute and continue with the next one
			cols[colIndex].attributes = append(cols[colIndex].attributes, att)
			continue
		}

		colIndex = len(cols)
		uniqCols[colIdent] = colIndex
		cols = append(cols, &column{ident: colIdent, attributes: []*data.Attribute{att}})

		if data.StoreCodecStdRecordValueJSONType.Is(att.Store) {
			// when dealing with encoded types there is probably
			// a different column that can handle the encoded payload
			cols[colIndex].columnType = &data.TypeJSON{}
			cols[colIndex].decode = decodeStdRecordValueJSON
			cols[colIndex].encode = encodeStdRecordValueJSON
		} else {
			cols[colIndex].columnType = att.Type
			if isSystemField(att.Ident) {
				cols[colIndex].decode = func(aa []*data.Attribute, raw any, r *types.Record) error {
					return decodeSystemField(aa[0].Ident, raw, r)
				}

				cols[colIndex].encode = func(aa []*data.Attribute, r *types.Record) (any, error) {
					// ignoring attribute slice because it's strictly one field we will be dealing with
					// + only interested  in the single-value (first) fields

					return getSystemFieldValue(r, aa[0].Ident)
				}

			} else {
				cols[colIndex].decode = decodeRecordValue
				cols[colIndex].encode = func(aa []*data.Attribute, r *types.Record) (any, error) {
					// ignoring attribute slice because it's strictly one field we will be dealing with
					// + only interested  in the single-value (first) fields
					//
					// any potential incompatibilities should be resolved before this point!
					if v := r.Values.Get(aa[0].Ident, 0); v != nil {
						return v.Value, nil
					}

					return nil, nil
				}

			}
		}
	}

	ms.columns = cols

	return nil
}

func (ms *crs) Truncate(ctx context.Context) error {
	sql, args, err := ms.truncateSql().ToSQL()
	if err != nil {
		return err
	}

	_, err = ms.conn.ExecContext(ctx, sql, args...)
	return err
}

func (ms *crs) Create(ctx context.Context, rr ...*types.Record) error {
	sql, args, err := ms.insertSql(rr...).ToSQL()
	if err != nil {
		return err
	}

	_, err = ms.conn.ExecContext(ctx, sql, args...)
	return err
}

func (ms *crs) Update(ctx context.Context, r *types.Record) error {
	sql, args, err := ms.updateSql(r).ToSQL()
	if err != nil {
		return err
	}

	_, err = ms.conn.ExecContext(ctx, sql, args...)
	return err
}

func (ms *crs) Delete(ctx context.Context, r *types.Record) error {
	sql, args, err := ms.deleteByIdSql(r.ID).ToSQL()
	if err != nil {
		return err
	}

	_, err = ms.conn.ExecContext(ctx, sql, args...)
	return err
}

func (ms *crs) Search(ctx context.Context, f types.RecordFilter) (i *iterator, err error) {
	// construct base query
	sql := ms.searchSql(f)
	_ = sql
	return &iterator{ms: ms}, nil
}

//func Search(ctx context.Context, s querier, m *data.Model, f types.RecordFilter) (set types.RecordSet, _ types.RecordFilter, err error) {
//	// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
//	f.PrevPage, f.NextPage = nil, nil
//
//	if f.PageCursor != nil {
//		// Page cursor exists; we need to validate it against used sort
//		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
//		// from the cursor.
//		// This (extracted sorting info) is then returned as part of response
//		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
//			return
//		}
//	}
//
//	// Make sure results are always sorted at least by primary keys
//	if f.Sort.Get("id") == nil {
//		f.Sort = append(f.Sort, &filter.SortExpr{
//			Column:     "id",
//			Descending: f.Sort.LastDescending(),
//		})
//	}
//
//	// Cloned sorting instructions for the actual sorting
//	// Original are passed to the etchFullPageOfApplications fn used for cursor creation;
//	// direction information it MUST keep the initial
//	sort := f.Sort.Clone()
//
//	// When cursor for a previous page is used it's marked as reversed
//	// This tells us to flip the descending flag on all used sort keys
//	if f.PageCursor != nil && f.PageCursor.ROrder {
//		sort.Reverse()
//	}
//
//	set, f.PrevPage, f.NextPage, err = fetchFullPage(ctx, s, m, f, sort)
//
//	f.PageCursor = nil
//	if err != nil {
//		return nil, f, err
//	}
//
//	return set, f, nil
//
//}
//
//func fetchFullPage(ctx context.Context, s querier, m *data.Model, filter types.RecordFilter, sort filter.SortExprSet) (set []*types.Record, prev, next *filter.PagingCursor, err error) {
//	var (
//		aux []*types.Record
//
//		// When cursor for a previous page is used it's marked as reversed
//		// This tells us to flip the descending flag on all used sort keys
//		reversedOrder = filter.PageCursor != nil && filter.PageCursor.ROrder
//
//		// Copy no. of required items to limit
//		// Limit will change when doing subsequent queries to fill
//		// the set with all required items
//		limit = filter.Limit
//
//		reqItems = filter.Limit
//
//		// cursor to prev. page is only calculated when cursor is used
//		hasPrev = filter.PageCursor != nil
//
//		// next cursor is calculated when there are more pages to come
//		hasNext bool
//
//		tryFilter types.RecordFilter
//	)
//
//	set = make([]*types.Record, 0, rdbms.DefaultSliceCapacity)
//
//	for try := 0; try < rdbms.MaxRefetches; try++ {
//		// Copy filter & apply custom sorting that might be affected by cursor
//		tryFilter = filter
//		tryFilter.Sort = sort
//
//		if limit > 0 {
//			// fetching + 1 to peak ahead if there are more items
//			// we can fetch (next-page cursor)
//			tryFilter.Limit = limit + 1
//		}
//
//		if aux, hasNext, err = query(ctx, s, m, tryFilter); err != nil {
//			return nil, nil, nil, err
//		}
//
//		if len(aux) == 0 {
//			// nothing fetched
//			break
//		}
//
//		// append fetched items
//		set = append(set, aux...)
//
//		if reqItems == 0 || !hasNext {
//			// no max requested items specified, break out
//			break
//		}
//
//		collected := uint(len(set))
//
//		if reqItems > collected {
//			// not enough items fetched, try again with adjusted limit
//			limit = reqItems - collected
//
//			if limit < rdbms.MinEnsureFetchLimit {
//				// In case limit is set very low and we've missed records in the first fetch,
//				// make sure next fetch limit is a bit higher
//				limit = rdbms.MinEnsureFetchLimit
//			}
//
//			// Update cursor so that it points to the last item fetched
//			tryFilter.PageCursor = collectCursorValues(set[collected-1], filter.Sort...)
//
//			// Copy reverse flag from sorting
//			tryFilter.PageCursor.LThen = filter.Sort.Reversed()
//			continue
//		}
//
//		if reqItems < collected {
//			set = set[:reqItems]
//		}
//
//		break
//	}
//
//	collected := len(set)
//
//	if collected == 0 {
//		return nil, nil, nil, nil
//	}
//
//	if reversedOrder {
//		// Fetched set needs to be reversed because we've forced a descending order to get the previous page
//		for i, j := 0, collected-1; i < j; i, j = i+1, j-1 {
//			set[i], set[j] = set[j], set[i]
//		}
//
//		// when in reverse-order rules on what cursor to return change
//		hasPrev, hasNext = hasNext, hasPrev
//	}
//
//	if hasPrev {
//		prev = collectCursorValues(set[0], filter.Sort...)
//		prev.ROrder = true
//		prev.LThen = !filter.Sort.Reversed()
//	}
//
//	if hasNext {
//		next = collectCursorValues(set[collected-1], filter.Sort...)
//		next.LThen = filter.Sort.Reversed()
//	}
//
//	return set, prev, next, nil
//}
//
//func collectCursorValues(res *types.Record, cc ...*filter.SortExpr) *filter.PagingCursor {
//	var (
//		cur = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}
//
//		hasUnique bool
//
//		pkID bool
//
//		collect = func(cc ...*filter.SortExpr) {
//			for _, c := range cc {
//				switch c.Column {
//				case sysID:
//					cur.Set(c.Column, res.ID, c.Descending)
//					pkID = true
//				case "createdAt":
//					cur.Set(c.Column, res.CreatedAt, c.Descending)
//				case "updatedAt":
//					cur.Set(c.Column, res.UpdatedAt, c.Descending)
//				case "deletedAt":
//					cur.Set(c.Column, res.DeletedAt, c.Descending)
//				}
//			}
//		}
//	)
//
//	collect(cc...)
//	if !hasUnique || !pkID {
//		collect(&filter.SortExpr{Column: "id", Descending: false})
//	}
//
//	return cur
//}
//
//func query(ctx context.Context, s querier, m *data.Model, f types.RecordFilter) (_ []*types.Record, more bool, err error) {
//	var (
//		ok bool
//
//		set         = make([]*types.Record, 0, rdbms.DefaultSliceCapacity)
//		res         *types.Record
//		rows        *sql.Rows
//		count       uint
//		expr, tExpr []goqu.Expression
//
//		sortExpr []exp.OrderedExpression
//
//		sql *goqu.SelectDataset
//	)
//
//	if err != nil {
//		err = fmt.Errorf("could generate filter expression for Application: %w", err)
//		return
//	}
//
//	// paging feature is enabled
//	if f.PageCursor != nil {
//		if tExpr, err = rdbms.Cursor(f.PageCursor); err != nil {
//			return
//		} else {
//			expr = append(expr, tExpr...)
//		}
//	}
//
//	sql, err = ms.searchSql(f)
//	sql.Where(expr...)
//
//	//query, := searchSql(m, f).Where(expr...)
//
//	// sorting feature is enabled
//	if sortExpr, err = rdbms.Order(f.Sort, nil); err != nil {
//		err = fmt.Errorf("could generate order expression for Application: %w", err)
//		return
//	}
//
//	if len(sortExpr) > 0 {
//		sql = sql.Order(sortExpr...)
//	}
//
//	if f.Limit > 0 {
//		sql = sql.Limit(f.Limit)
//	}
//
//	rows, err = s.Query(ctx, sql)
//	if err != nil {
//		err = fmt.Errorf("could not query Application: %w", err)
//		return
//	}
//
//	if err = rows.Err(); err != nil {
//		err = fmt.Errorf("could not query Application: %w", err)
//		return
//	}
//
//	defer func() {
//		closeError := rows.Close()
//		if err == nil {
//			// return error from close
//			err = closeError
//		}
//	}()
//
//	for rows.Next() {
//		if err = rows.Err(); err != nil {
//			err = fmt.Errorf("could not query Application: %w", err)
//			return
//		}
//
//		listOfSelectedColumns := make([]any, 10)
//
//		// @todo this does not work!
//
//		if err = rows.Scan(listOfSelectedColumns...); err != nil {
//			err = fmt.Errorf("could not scan rows for Application: %w", err)
//			return
//		}
//
//		count++
//		// @todo convert listOfSelectedCOlumns into *types.Record
//
//		//if res, err = aux.decode(); err != nil {
//		//	err = fmt.Errorf("could not decode Application: %w", err)
//		//	return
//		//}
//
//		// check fn set, call it and see if it passed the test
//		// if not, skip the item
//		if f.Check != nil {
//			if ok, err = f.Check(res); err != nil {
//				return
//			} else if !ok {
//				continue
//			}
//		}
//
//		set = append(set, res)
//	}
//
//	return set, f.Limit > 0 && count >= f.Limit, err
//
//}

// constructs SQL for selecting records from a table, converting parts of record filter into conditions
func (ms *crs) searchSql(f types.RecordFilter) *goqu.SelectDataset {
	var (
		err  error
		base = ms.selectSql()
		tmp  exp.Expression
		cnd  []exp.Expression
	)

	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return base.SetError(err)
		}
	}

	if f.Sort.Get(sysID) == nil {
		// Make sure results are always sorted at least by primary key
		f.Sort = append(f.Sort, &filter.SortExpr{
			Column:     sysID,
			Descending: f.Sort.LastDescending(),
		})
	}

	{
		// Add module & namespace constraints when model expects (has configured attributes) for them
		//
		// This covers both scenarios:
		// 1) Model is configured to store records in a dedicated table
		//    without module and/or namespace attributes
		//
		// 2) Model has module and/or namespace attribute and saves records
		//    from different modules in the same table

		if ms.sysExprNamespaceID != nil {
			cnd = append(cnd, ms.sysExprNamespaceID.Eq(f.NamespaceID))
		} else {
			// @todo check if f.NamespaceID is compatible
		}

		if ms.sysExprModuleID != nil {
			cnd = append(cnd, ms.sysExprModuleID.Eq(f.ModuleID))
		} else {
			// @todo check if f.ModuleID is compatible
		}
	}

	{
		if len(f.LabeledIDs) > 0 {
			// Limit by LabeledIDs (list of record IDs)
			cnd = append(cnd, ms.sysColumnID.In(f.LabeledIDs))
		}
	}

	{
		// If module supports soft-deletion (= delete-at attribute is present)
		// we need to make sure we respect it
		if ms.sysExprDeletedAt != nil {
			switch f.Deleted {
			case filter.StateExclusive:
				// only not-null values
				cnd = append(cnd, ms.sysExprDeletedAt.IsNotNull())

			case filter.StateExcluded:
				// exclude all non-null values
				cnd = append(cnd, ms.sysExprDeletedAt.IsNull())
			}
		}
	}

	if len(strings.TrimSpace(f.Query)) > 0 {
		if tmp, err = ms.queryParser.Parse(f.Query); err != nil {
			return base.SetError(err)
		}

		cnd = append(cnd, tmp)
	}

	return ms.selectSql().Where(cnd...)
}

func (ms *crs) lookupByIdSql(id uint) *goqu.SelectDataset {
	return ms.selectSql().
		Where(ms.sysColumnID.Eq(id)).
		Limit(1)
}

func (ms *crs) selectSql() *goqu.SelectDataset {
	// working around a bug inside goqu lib that adds
	// * to the list of columns to be selected
	// even if we clear the columns first
	q := goqu.From(ms.table).Select(ms.table.Col(ms.columns[0].ident))

	for _, col := range ms.columns[1:] {
		q = q.SelectAppend(ms.table.Col(col.ident))
	}

	return q
}

func (ms *crs) truncateSql() (_ *goqu.TruncateDataset) {
	return goqu.Truncate(ms.table)
}

func (ms *crs) insertSql(rr ...*types.Record) (_ *goqu.InsertDataset) {
	var (
		err  error
		rows = make([]any, len(rr))
	)

	for i, r := range rr {
		rows[i], err = encode(ms.columns, r)
		if err != nil {
			return (&goqu.InsertDataset{}).SetError(err)
		}
	}

	return goqu.
		Insert(ms.table).
		Rows(rows...)
}

// updateSql generates SQL command for updating record
//
// Integrity check (ie module, namespace, dates, owners change) is out of scope
// for this layer
func (ms *crs) updateSql(r *types.Record) *goqu.UpdateDataset {
	encValues, err := encode(ms.columns, r, sysID)
	if err != nil {
		return (&goqu.UpdateDataset{}).SetError(err)
	}

	return goqu.
		Update(ms.table).
		Set(encValues).
		Where(ms.sysColumnID.Eq(r.ID))
}

func (ms *crs) deleteByIdSql(id uint64) *goqu.DeleteDataset {
	return goqu.Delete(ms.table).Where(ms.sysColumnID.Eq(id))
}

func (ms *crs) attrToExpr(ident string) (exp.LiteralExpression, error) {
	if !ms.model.HasAttribute(ident) {
		return nil, fmt.Errorf("unknown attribute %q", ident)
	}

	attr := ms.model.Attributes.FindByIdent(ident)
	switch s := attr.Store.(type) {
	case *data.StoreCodecAlias:
		// using column directly
		return exp.NewLiteralExpression(fmt.Sprintf("%q.%q", ms.model.Ident, s.Ident)), nil

	case *data.StoreCodecStdRecordValueJSON:
		// using JSON to handle embedded values
		return ms.deepJsonFn(s.Ident, attr.Ident, 0), nil
	}

	return exp.NewLiteralExpression(fmt.Sprintf("%q.%q", ms.model.Ident, ident)), nil
}

func attrColumnIdent(att *data.Attribute) string {
	switch ss := att.Store.(type) {
	case *data.StoreCodecStdRecordValueJSON:
		return ss.Ident

	case *data.StoreCodecAlias:
		return ss.Ident

	default:
		return att.Ident
	}
}

//func extractColumns(m *data.Model) []exp.IdentifierExpression {
//	var (
//		col  string
//		uniq = make(map[string]bool)
//		ii   = make([]exp.IdentifierExpression, 0, len(aa))
//	)
//
//	for _, a := range m.Attributes {
//		col = attrColumnIdent(a)
//		if uniq[col] {
//			continue
//		}
//
//		uniq[col] = true
//		ii = append(ii, col)
//	}
//
//	return
//}

// attributeExp converts attribute to sql expression
//func attributeExp(att *data.Attribute) (exp.Expression, error) {
//	if att.Store == nil {
//		return exp.NewIdentifierExpression("", "", att.Ident), nil
//
//	}
//	switch ss := att.Store.(type) {
//	case data.StoreCodecStdRecordValueJSON:
//		return exp.NewIdentifierExpression("", "", ss.Ident), nil
//	case data.StoreCodecAlias:
//		return exp.NewIdentifierExpression("", "", ss.Ident).As(att.Ident), nil
//
//	default:
//		return nil, fmt.Errorf("unknown store strategy %T for attribute %q", ss, att.Ident)
//	}
//}
