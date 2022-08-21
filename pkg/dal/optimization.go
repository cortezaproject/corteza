package dal

type (
	// opCost provides a general idea of expensive an operation is for a
	// specific pipeline step.
	//
	// dsSize provides a general idea of how large an underlaying dataset is.
	opCost       int
	dsSize       int
	stepAnalysis struct {
		scanCost   opCost
		searchCost opCost
		filterCost opCost
		sortCost   opCost

		outputSize dsSize
	}
)

var (
	pipelineOptimizers = []func(PipelineStep, bool) (PipelineStep, error){
		// @todo add more optimizers
	}
)

const (
	// operation computation indicators
	costUnknown opCost = iota
	costFree
	costCheep
	costAcceptable
	costExpensive
	costInfinite
)

const (
	// dataset size indicators
	sizeUnknown dsSize = iota
	sizeTiny
	sizeSmall
	sizeMedium
	sizeLarge
)
