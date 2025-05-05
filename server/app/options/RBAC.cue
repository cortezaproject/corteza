package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

RBAC: schema.#optionsGroup & {
	handle: "rbac"
	title:  "RBAC options"

	options: {
		log: {
			type:        "bool"
			description: "Log RBAC related events and actions"
		}

		max_index_size: {
			type:         "int"
			defaultGoExpr: "-1"
			description: """
				Limit the number of resources to keep in the in-memory index.
				When set to -1, max size is used.
				When set to 0, the in memory index is not used.
				"""
		}

		synchronous: {
			type:         "bool"
			defaultGoExpr: "false"
			description: """
				Synchronous lets us make all the procedures synchronous for ease of testing
				This should always be false in production
				"""
		}

		reindex_strategy: {
			type:         "string"
			defaultValue: ""
			description: """
				Reindex strategy defines what strategy we should use.
				The available options are:

				. `memory`: prioritize memory consumption which reduces performance during reindexing.
				. `speed`: prioritize speed during reindexing; memory consumption will be 2n where n is the current index size.

				If you wish to prioritize memory and speed, consider using `speed` with a lower max index size
				"""
		}

		decay_factor: {
			type:         "float64"
			defaultGoExpr: "0.9"
			description: """
				Decay factor controls how long an item should be kept in the index while not in use.
				"""
		}

		decay_interval: {
			type:          "time.Duration"
			defaultGoExpr: "time.Minute * 30"
			description: """
				Decay interval controls how fast the decay factor is applied to the index key.
				"""
		}

		cleanup_interval: {
			type:          "time.Duration"
			defaultGoExpr: "time.Minute * 31"
			description: """
				Cleanup interval controls when unused/low-scored index items should be yanked out of the index counter.
				"""
		}

		reindex_interval: {
			type:          "time.Duration"
			defaultGoExpr: "time.Minute * 10"
			description: """
				Reindex interval controls when the index should be re-calculated.
				"""
		}

		index_flush_interval: {
			type:          "time.Duration"
			defaultGoExpr: "time.Minute * 35"
			description: """
				[IMPORTANT]
				====
				Unused, will be added when state preservation is implemented.
				====
				"""
		}

		service_user: {}
		bypass_roles: {
			defaultValue: "super-admin"
			description: """
				Space delimited list of role handles.
				These roles causes short-circuiting access control check and allowing all operations.
				System will refuse to start if check-bypassing roles are also listed as authenticated or anonymous auto-assigned roles.
				"""
		}
		authenticated_roles: {
			defaultValue: "authenticated"
			description: """
				Space delimited list of role handles.
				These roles are automatically assigned to authenticated user.
				Memberships can not be managed for these roles.
				System will refuse to start if roles listed here are also listed under anonymous roles
				"""
		}
		anonymous_roles: {
			defaultValue: "anonymous"
			description: """
				Space delimited list of role handles.
				These roles are automatically assigned to anonymous user.
				Memberships can not be managed for these roles.
				"""
		}
	}
}
