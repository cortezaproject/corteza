package provision

import (
	"context"
	"fmt"
	"net/http"

	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
)

func apigwFilters(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	log.Info("provision start")
	defer log.Info("provision end")

	var (
		filters types.ApigwFilterSet
	)

	if filters, _, err = store.SearchApigwFilters(ctx, s, types.ApigwFilterFilter{Ref: "jsonResponse"}); err != nil {
		return
	}

	for _, f := range filters {
		h := http.Header{}
		h.Add("Content-Type", "application/json")

		f.Ref = "response"
		f.Params = map[string]interface{}{
			"header": h,
			"input":  &atypes.Expr{Expr: f.Params["input"].(string), Type: "String"},
		}

		if err = store.UpdateApigwFilter(ctx, s, f); err != nil {
			log.Warn(fmt.Sprintf("could not migrate jsonResponse to response: %s", err))
		}
	}

	return
}
