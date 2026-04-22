package schema

#platform: {
	ident: #baseHandle | *"corteza"

	options: [...#optionsGroup]

	components: [...{platform: ident} & #component]

	resources: {
		[key=#handle]: #Resource & {
			"handle": key,
			"platform": ident
		}
	}

	// automation: {
	//  types: ....
	//  function ....
	// }
}
