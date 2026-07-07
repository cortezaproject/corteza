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

	// bundles: standalone type registries for non-component packages
	bundles: [...#TypeBundle] | *[]

	// automation: {
	//  types: ....
	//  function ....
	// }
}
