package postgres

//var (
//	sqlExprRegistry = map[string]rdbms.HandlerSig{
//		// functions
//		// - filtering
//		"quarter": makeGenericExtrFncHandler("QUARTER"),
//		"year":    makeGenericExtrFncHandler("YEAR"),
//		"month":   makeGenericExtrFncHandler("MONTH"),
//		"date":    makeGenericExtrFncHandler("DAY"),
//	}
//)
//
//func sqlASTFormatter(n *ql.ASTNode) rdbms.HandlerSig {
//	return sqlExprRegistry[n.Ref]
//}
//
//func makeGenericExtrFncHandler(extr string) rdbms.HandlerSig {
//	return func(aa ...rdbms.FormattedASTArgs) (out string, args []interface{}, selfEnclosed bool, err error) {
//		if len(aa) != 1 {
//			err = fmt.Errorf("expecting 1 arguments, got %d", len(aa))
//			return
//		}
//
//		out = fmt.Sprintf("EXTRACT(%s FROM %s)", extr, aa[0].S)
//		args = aa[0].Args
//		return
//	}
//}
