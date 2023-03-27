package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

attachment: schema.#optionsGroup & {
	handle: "attachment"

	options: {
		avatar_max_file_size: {
			type: "int64"
			defaultGoExpr: "1000000"
			description:  "Avatar image maximum upload size, default value is 1MB"
		}
		avatar_initials_font_path: {
			defaultValue: "fonts/Poppins-Regular.ttf"
			description:  "Avatar initials font file path"
			env:          "AVATAR_INITIALS_FONT_PATH"
		}
		avatar_initials_background_color: {
			defaultValue: "#F3F3F3"
			description:  "Avatar initials background color"
			env:          "AVATAR_INITIALS_BACKGROUND_COLOR"
		}
		avatar_initials_color: {
			defaultValue: "#162425"
			description:  "Avatar initials text color"
			env:          "AVATAR_INITIALS_COLOR"
		}
	}
}
