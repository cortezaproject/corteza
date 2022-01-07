package schema

#platform: {
	ident: #baseHandle | *"corteza"

	components: [...{platform: ident} & #component]
	gig: {
		decoders: [ ...({kind: "decoder"} & #gigTask)]
		preprocessors: [ ...({kind: "preprocessor"} & #gigTask)]
		postprocessors: [ ...({kind: "postprocessor"} & #gigTask)]

		workers: [...#gigWorker]
	}

	// env-var definitions
	// options: {}

	//

	// automation: {
	//  types: ....
	//  function ....
	// }
}
