package schema

#platform: {
	ident: #baseHandle | *"human"

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
