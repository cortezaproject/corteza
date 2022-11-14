package schema

#platform: {
	ident: #baseHandle | *"corteza"

	options: [...#optionsGroup]
	components: [...{platform: ident} & #component]
	// env-var definitions
	// options: {}

	//

	// automation: {
	//  types: ....
	//  function ....
	// }
}
