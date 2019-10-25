package types

type (
	Settings struct {
		General struct {
			Mail struct {
				Logo   string
				Header string `kv:"header.en"`
				Footer string `kv:"footer.en"`
			}
		}
	}
)
