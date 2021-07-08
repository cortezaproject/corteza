package commands

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cobra"
)

func RBAC(ctx context.Context, storeInit func(ctx context.Context) (store.Storer, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rbac",
		Short: "RBAC tools",
		Long:  "Check and manipulates permissions",
	}

	cmd.AddCommand(rbacList(ctx, storeInit))

	// @todo command that can grant/revoke/reset-all permissions
	//       in a similar format(s) as we do listing so that users can
	//       copy-paste the whole output, modify it and import it back

	return cmd
}

func rbacList(ctx context.Context, storeInit func(ctx context.Context) (store.Storer, error)) (cmd *cobra.Command) {
	var (
		resources   []string
		roles       []string
		operations  []string
		groupBy     string
		allow, deny bool

		matchResources = func(r *rbac.Rule) bool {
			if len(resources) == 0 {
				return true
			}

			for _, res := range resources {
				if strings.HasPrefix(r.Resource, res) {
					return true
				}
			}

			return false
		}

		// makes a simple utiliy function for matching rules according to the used flags
		ruleMatcher = func(s store.Storer) (map[uint64]*types.Role, func(r *rbac.Rule) bool) {
			var (
				opsMap    = slice.ToStringBoolMap(operations)
				rr        []*types.Role
				rolAuxMap = slice.ToStringBoolMap(roles)
				rolMap    = make(map[uint64]*types.Role)
				rolMatch  = make(map[uint64]bool)
				err       error
			)

			rr, _, err = store.SearchRoles(ctx, s, types.RoleFilter{})
			for _, r := range rr {
				rolMap[r.ID] = r
				if rolAuxMap[r.Name] || rolAuxMap[r.Handle] || rolAuxMap[strconv.FormatUint(r.ID, 10)] {
					rolMatch[r.ID] = true
				}
			}

			cli.HandleError(err)

			return rolMap, func(r *rbac.Rule) bool {
				if !matchResources(r) {
					return false
				}
				if len(opsMap) > 0 && !opsMap[r.Operation] {
					return false
				}

				// this use of rolAuxMap to check if there are roles specified in the flags
				// and rolMap for checking if for actual role map is intentional!
				if len(rolAuxMap) > 0 && !rolMatch[r.RoleID] {
					return false
				}

				if allow && r.Access != rbac.Allow {
					return false
				}

				if deny && r.Access != rbac.Deny {
					return false
				}

				return true
			}
		}

		ruleSorter = func(rr []*rbac.Rule) {
			sort.SliceStable(rr, func(i, j int) bool {
				switch groupBy {
				case "role", "rMap":
					return rr[i].RoleID < rr[j].RoleID
				case "res", "resource", "resources":
					return strings.Compare(rr[i].Resource, rr[i].Resource) < 0
				case "op", "ops", "operation", "operations":
					return strings.Compare(rr[i].Operation, rr[j].Operation) < 0
				}

				return rr[i].Access < rr[j].Access
			})
		}

		roleDisplayName = func(rMap map[uint64]*types.Role, r *rbac.Rule) (role string) {
			if rMap[r.RoleID] != nil {
				role = rMap[r.RoleID].Name
				if role == "" {
					role = rMap[r.RoleID].Handle
				}
			}

			if role == "" {
				return strconv.FormatUint(r.RoleID, 10)
			}

			return
		}

		lengths = func(rr []*rbac.Rule, rMap map[uint64]*types.Role) (op, res, role int) {
			for _, r := range rr {
				if len(r.Operation) > op {
					op = len(r.Operation)
				}
				if len(r.Resource) > res {
					res = len(r.Resource)
				}
			}
			for _, r := range rMap {
				if len(r.Name) > role {
					role = len(r.Name)
				}
				if len(r.Handle) > role {
					role = len(r.Handle)
				}
			}

			return
		}
	)

	cmd = &cobra.Command{
		Use:   "list",
		Short: "Check applied permissions against given file (only supports compose permissions for now)",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				rr      []*rbac.Rule
				s, err  = storeInit(ctx)
				role    string
				gBucket string
			)

			cli.HandleError(err)

			rr, _, err = store.SearchRbacRules(cli.Context(), s, rbac.RuleFilter{})
			cli.HandleError(err)

			ruleSorter(rr)
			rMap, matchRule := ruleMatcher(s)

			longestOp, longestRes, longestRole := lengths(rr, rMap)

			hr := strings.Repeat("-", longestOp+longestRes+longestRole+20)

			for i, r := range rr {
				if !matchRule(r) {
					continue
				}

				r.Operation = r.Operation + strings.Repeat(" ", longestOp-len(r.Operation))
				r.Resource = r.Resource + strings.Repeat(" ", longestRes-len(r.Resource))

				role = roleDisplayName(rMap, r)
				role = role + strings.Repeat(" ", longestRole-len(role))

				bOp := strings.Repeat(" ", longestOp)
				bRes := strings.Repeat(" ", longestRes)
				bRole := strings.Repeat(" ", longestRole)

				switch groupBy {
				case "role", "roles":
					if gBucket != role {
						if i != 0 {
							cmd.Println(hr)
						}
						cmd.Printf("%s %7s %s on %s\n", role, r.Access, r.Operation, r.Resource)
					} else {
						cmd.Printf("%s %7s %s on %s\n", bRole, r.Access, r.Operation, r.Resource)
					}

					gBucket = role

				case "res", "resource", "resources":
					if gBucket != r.Resource {
						if i != 0 {
							cmd.Println(hr)
						}

						cmd.Printf("on %s %7s %s to %s\n", r.Resource, r.Access, role, r.Operation)
					} else {
						cmd.Printf("   %s %7s %s to %s\n", bRes, r.Access, role, r.Operation)
					}

					gBucket = r.Resource

				case "op", "ops", "operation", "operations":
					if gBucket != r.Operation {
						if i != 0 {
							cmd.Println(hr)
						}

						cmd.Printf("%s %7s %s to %s\n", r.Operation, r.Access, role, r.Resource)
					} else {
						cmd.Printf("%s %7s %s to %s\n", bOp, r.Access, role, r.Resource)
					}

					gBucket = r.Operation
				default:
					cmd.Printf("%7s %s to %s on %s\n", r.Access, role, r.Operation, r.Resource)
				}
			}

		},
	}

	cmd.Flags().StringArrayVarP(&resources, "resource", "r", nil, "Filter by resource (by prefix)")
	cmd.Flags().StringArrayVarP(&roles, "role", "", nil, "Filter by role (handle or ID)")
	cmd.Flags().StringArrayVarP(&operations, "operation", "o", nil, "Filter by operation")
	cmd.Flags().BoolVarP(&deny, "deny", "d", false, "Show only deny")
	cmd.Flags().BoolVarP(&allow, "allow", "a", false, "Show only allows")
	cmd.Flags().StringVarP(&groupBy, "group", "g", "", "Group rules on output")

	return
}
