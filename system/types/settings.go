package types

type (
	Settings struct {
		General struct {
			Mail struct {
				Logo   string
				Header string `kv:"footer.en"`
				Footer string `kv:"header.en"`
			}
		}
	}
)
