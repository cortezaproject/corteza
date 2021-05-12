package expr

func Example_toJSON() {
	var (
		p = map[string]interface{}{
			"vars": Must(NewVars(
				&Vars{value: map[string]TypedValue{
					"k1": &String{value: "v1"},
					"k2": &String{value: "v2"},
				}},
			)),
		}
	)

	eval(`toJSON(vars)`, p)

	// output:
	// {"k1":{"@value":"v1","@type":"String"},"k2":{"@value":"v2","@type":"String"}}
}
