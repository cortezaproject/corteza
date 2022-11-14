package indexer

var (
	status struct {
		Index     string `json:"index"`
		Timestamp string `json:"timestamp"`
		Full      bool
	}
)

//func pushStatus(ctx context.Context, esc *elasticsearch.Client) error {
//	var (
//		buf = &bytes.Buffer{}
//		st = struct {
//
//		}
//	)
//
//	json.NewEncoder(buf).Encode()
//
//	spew.Dump(esc.Create("index", "", ))
//
//	return nil
//}
