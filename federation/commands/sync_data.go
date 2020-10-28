package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	ct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/federation/service"
	"github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/decoder"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var (
	dataSyncer *service.Syncer
)

func commandSyncData(ctx context.Context) func(*cobra.Command, []string) {

	// from
	//   - module id (Account) = 167296235227578369
	//   - namespace id = 167296235059806209

	// to
	//   - module id (Account) = 167296265594339329
	//   - namespace id = 167296264570929153

	return func(_ *cobra.Command, _ []string) {

		// todo
		//  - get all nodes and loop
		//  - get all shared modules that have a field mapping
		//  - for each shared module, the add method needs shared module info as payload

		// todo - get node by id
		const (
			nodeID    = 276342359342989444
			moduleID  = 376342359342989444
			limit     = 5
			syncDelay = 10
		)

		node, err := service.DefaultNode.FindByID(ctx, nodeID)

		if err != nil {
			service.DefaultLogger.Info("could not find any nodes to sync")
			return
		}

		ticker := time.NewTicker(time.Second * syncDelay)
		basePath := fmt.Sprintf("/federation/nodes/%d/modules/%d/records/", node.SharedNodeID, moduleID)

		dataSyncer = service.NewSyncer()

		// todo - get auth from the node
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// get the last sync per-node
		lastSync := getLastSyncTime(ctx, nodeID, types.NodeSyncTypeData)

		ex := ""
		if lastSync != nil {
			ex = fmt.Sprintf(" from last sync on (%s)", lastSync.Format(time.RFC3339))
		}

		service.DefaultLogger.Info(fmt.Sprintf("Starting structure sync%s", ex))

		url := types.SyncerURI{
			BaseURL:  baseURL,
			Path:     basePath,
			Limit:    limit,
			LastSync: lastSync,
		}

		go queueUrl(&url, urls, dataSyncer)

		for {
			select {
			case <-ctx.Done():
				// do any cleanup here
				service.DefaultLogger.Info(fmt.Sprintf("Stopping sync [processed: %d, persisted: %d]", countProcess, countPersist))
				return
			case <-ticker.C:
				lastSync := getLastSyncTime(ctx, nodeID, types.NodeSyncTypeStructure)

				url := types.SyncerURI{
					BaseURL:  baseURL,
					Path:     basePath,
					Limit:    limit,
					LastSync: lastSync,
				}

				if lastSync == nil {
					service.DefaultLogger.Info(fmt.Sprintf("Start fetching modules from beginning of time, since lastSync == nil"))
				} else {
					service.DefaultLogger.Info(fmt.Sprintf("Start fetching modules, start with time: %s", lastSync.Format(time.RFC3339)))
				}

				go queueUrl(&url, urls, dataSyncer)
			case url := <-urls:
				select {
				case <-ctx.Done():
					// do any cleanup here
					service.DefaultLogger.Info(fmt.Sprintf("Stopping sync [processed: %d, persisted: %d]", countProcess, countPersist))
					return
				default:
				}

				s, err := url.String()

				if err != nil {
					continue
				}

				responseBody, err := dataSyncer.Fetch(ctx, s)

				if err != nil {
					continue
				}

				payloads <- responseBody
			case p := <-payloads:
				body, err := ioutil.ReadAll(p)

				// handle error
				if err != nil {
					continue
				}

				u := types.SyncerURI{
					BaseURL: baseURL,
					Path:    basePath,
					Limit:   limit,
				}

				// method needs:
				//  - compose module ID (from module_mapping)
				//  - namespace ID (from module mapping)

				err = dataSyncer.Process(ctx, body, node.ID, urls, u, dataSyncer, persistExposedRecordData)

				if err != nil {
					// handle error
					service.DefaultLogger.Error(fmt.Sprintf("error on handling payload: %s", err))
				} else {
					n := time.Now().UTC()

					// add to db - nodes_sync
					new := &types.NodeSync{
						NodeID:       node.ID,
						SyncStatus:   types.NodeSyncStatusSuccess,
						SyncType:     types.NodeSyncTypeStructure,
						TimeOfAction: n,
					}

					new, err := service.DefaultNodeSync.Create(ctx, new)

					if err != nil {
						service.DefaultLogger.Warn(fmt.Sprintf("could not update sync status: %s", err))
					}
				}
			}
		}
	}
}

// persistExposedRecordData gets the payload from syncer and
// uses the decode package to decode the whole set, depending on
// the filtering that was used (limit)
func persistExposedRecordData(ctx context.Context, payload []byte, nodeID uint64, dataSyncer *service.Syncer) error {
	countProcess = countProcess + 1

	// now := time.Now()
	o, err := decoder.DecodeFederationRecordSync([]byte(payload))

	if err != nil {
		return err
	}

	if len(o) == 0 {
		return nil
	}

	service.DefaultLogger.Info(fmt.Sprintf("Adding %d objects", len(o)))

	for _, er := range o {
		// spew.Dump(o)
		// create a new compose record
		// get the compose module for this record
		// fill in the values

		// tmp
		vals := []*ct.RecordValue{
			&ct.RecordValue{
				Name:  "AccountName",
				Value: "Accout name is required",
			},
		}

		for _, v := range er.Values {
			if v.Name == "AccountName" {
				vv := &ct.RecordValue{
					Name:  "Name",
					Value: v.Value,
				}
				vals = append(vals, vv)
			}

			if v.Name == "Facebook" {
				vv := &ct.RecordValue{
					Name:  "Facebook",
					Value: v.Value,
				}
				vals = append(vals, vv)
			}

			if v.Name == "Phone" {
				vv := &ct.RecordValue{
					Name:  "Phone",
					Value: v.Value,
				}
				vals = append(vals, vv)
			}
		}

		//   - module id (Account) = 167296265594339329
		//   - namespace id = 167296264570929153

		rec := &ct.Record{
			ModuleID:    167296265594339329,
			NamespaceID: 167296264570929153,
			Values:      vals,
		}

		rc, err := dataSyncer.CService.With(ctx).Create(rec)

		spew.Dump(rc, err)

	}

	return nil
}
