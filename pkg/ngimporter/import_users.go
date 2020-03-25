package ngimporter

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"time"

	"github.com/cortezaproject/corteza-server/compose/repository"
	cct "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/ngimporter/types"
	sysRepo "github.com/cortezaproject/corteza-server/system/repository"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

// imports system users based on the provided source
func importUsers(ctx context.Context, is *types.ImportSource, ns *cct.Namespace) (map[string]uint64, *types.ImportSource, error) {
	db := repository.DB(ctx)
	repoUser := sysRepo.User(ctx, db)
	// this provides a map between importSourceID -> CortezaID
	mapping := make(map[string]uint64)

	// create a new buffer for user object, so we don't loose our data
	var bb bytes.Buffer
	ww := csv.NewWriter(&bb)
	defer ww.Flush()

	// get fields
	r := csv.NewReader(is.Source)
	header, err := r.Read()
	if err != nil {
		return nil, nil, err
	}
	ww.Write(header)

	// create users
	for {
	looper:
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, nil, err
		}

		ww.Write(record)

		u := &sysTypes.User{}
		for i, h := range header {
			val := record[i]

			// when creating users we only care about a handfull of values.
			// the rest are included in the module
			switch h {
			case "Username":
				u.Username = record[i]
				break

			case "Email":
				u.Email = record[i]
				break

			case "FirstName":
				u.Name = record[i]
				break

			case "LastName":
				u.Name = u.Name + " " + record[i]
				break

			case "Alias":
				u.Handle = record[i]
				break

			case "CreatedDate":
				if val != "" {
					u.CreatedAt, err = time.Parse(types.SfDateTimeLayout, val)
					if err != nil {
						return nil, nil, err
					}
				}
				break

			case "LastModifiedDate":
				if val != "" {
					tt, err := time.Parse(types.SfDateTimeLayout, val)
					u.UpdatedAt = &tt
					if err != nil {
						return nil, nil, err
					}
				}
				break

				// ignore deleted values, as SF provides minimal info about those
			case "IsDeleted":
				if val == "1" {
					goto looper
				}
			}
		}

		// this allows us to reuse existing users
		uu, err := repoUser.FindByEmail(u.Email)
		if err == nil {
			u = uu
		} else {
			u, err = repoUser.Create(u)
			if err != nil {
				return nil, nil, err
			}
		}

		mapping[record[0]] = u.ID
	}

	nis := &types.ImportSource{
		Name:    is.Name,
		Header:  is.Header,
		DataMap: is.DataMap,
		Path:    is.Path,
		Source:  &bb,
	}

	return mapping, nis, nil
}
