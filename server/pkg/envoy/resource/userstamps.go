package resource

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	Userstamp struct {
		UserID uint64
		Ref    string
		U      *types.User
		// S is the stringified user representation
		S string
	}

	Userstamps struct {
		CreatedBy *Userstamp
		UpdatedBy *Userstamp
		DeletedBy *Userstamp
		OwnedBy   *Userstamp
		RunAs     *Userstamp
	}
	UserstampIndex map[string]*Userstamp
)

func (us *Userstamp) MarshalYAML() (interface{}, error) {
	if us == nil {
		return nil, nil
	}

	if us.U != nil {
		if us.U.Handle != "" {
			return us.U.Handle, nil
		}
		if us.U.Username != "" {
			return us.U.Username, nil
		}
		if us.U.Email != "" {
			return us.U.Email, nil
		}
		if us.U.Name != "" {
			return us.U.Name, nil
		}
	}

	if us.Ref != "" {
		return us.Ref, nil
	}

	if us.UserID > 0 {
		return us.UserID, nil
	}

	return nil, errors.New("invalid userstamp")
}

func (us *Userstamp) MarshalJSON() ([]byte, error) {
	if us == nil {
		return nil, nil
	}

	l := ""

	if us.U != nil {
		if us.U.Handle != "" {
			l = us.U.Handle
		}
		if us.U.Username != "" {
			l = us.U.Username
		}
		if us.U.Email != "" {
			l = us.U.Email
		}
		if us.U.Name != "" {
			l = us.U.Name
		}
	} else {
		if us.Ref != "" {
			l = us.Ref
		}

		if us.UserID > 0 {
			l = strconv.FormatUint(us.UserID, 10)
		}
	}

	if l == "" {
		return nil, errors.New("invalid userstamp")
	}

	return json.Marshal(l)
}

// MakeUserstamp initializes a userstamp from the passed user struct
func MakeUserstamp(u *types.User) *Userstamp {
	sID := strconv.FormatUint(u.ID, 10)
	return &Userstamp{
		UserID: u.ID,
		U:      u,
		Ref:    firstOkString(u.Handle, u.Email, u.Username, sID),
	}
}

// MakeUserstampFromRef initializes a userstamp from the passed reference
//
// If possible, the function determines the userID
func MakeUserstampFromRef(ref string) *Userstamp {
	id, err := strconv.ParseUint(ref, 10, 64)

	us := &Userstamp{}

	if err == nil && id != 0 {
		us.UserID = id
		us.U = &types.User{ID: id}
	}
	us.Ref = ref

	return us
}

// MakeUserstampFromRef initializes a userstamp from the passed userID
func MakeUserstampFromID(ID uint64) *Userstamp {
	if ID == 0 {
		return nil
	}
	return MakeUserstampFromRef(strconv.FormatUint(ID, 10))
}

// MakeUserstampsCUDO initializes userstamps for createdAt, updatedAt, deletedAt and ownedBy
func MakeUserstampsCUDO(c, u, d, o uint64) *Userstamps {
	us := &Userstamps{}

	if c > 0 {
		us.CreatedBy = &Userstamp{
			UserID: c,
			Ref:    strconv.FormatUint(c, 10),
		}
	}
	if u > 0 {
		us.UpdatedBy = &Userstamp{
			UserID: u,
			Ref:    strconv.FormatUint(u, 10),
		}
	}
	if d > 0 {
		us.DeletedBy = &Userstamp{
			UserID: d,
			Ref:    strconv.FormatUint(d, 10),
		}
	}
	if o > 0 {
		us.OwnedBy = &Userstamp{
			UserID: o,
			Ref:    strconv.FormatUint(o, 10),
		}
	}

	return us
}

// Model stringifies the userstamp
//
// @todo make this configurable
func (us *Userstamp) Model() (string, error) {
	if us == nil {
		return "", nil
	}

	if us.U != nil {
		if us.U.Handle != "" {
			return us.U.Handle, nil
		}
		if us.U.Username != "" {
			return us.U.Username, nil
		}
		if us.U.Email != "" {
			return us.U.Email, nil
		}
		if us.U.Name != "" {
			return us.U.Name, nil
		}
	}

	if us.Ref != "" {
		return us.Ref, nil
	}

	if us.UserID > 0 {
		return strconv.FormatUint(us.UserID, 10), nil
	}

	return "", errors.New("invalid userstamp")
}

func (ux UserstampIndex) Add(uu ...*types.User) {
	for _, u := range uu {
		sID := strconv.FormatUint(u.ID, 10)
		s := MakeUserstamp(u)

		ux[sID] = s
		ux[u.Email] = s
		if u.Handle != "" {
			ux[u.Handle] = s
		}
		if u.Username != "" {
			ux[u.Username] = s
		}
		if u.Name != "" {
			ux[u.Name+" "+u.Email] = s
		}
	}
}

func (ux UserstampIndex) GetByKey(kr interface{}) *Userstamp {
	if k, ok := kr.(string); ok {
		return ux[k]
	} else if k, ok := kr.(uint64); ok {
		return ux[strconv.FormatUint(k, 10)]
	}
	return nil
}

func (ux UserstampIndex) GetByStamp(s *Userstamp) *Userstamp {
	if s == nil {
		return nil
	}

	if s.Ref != "" {
		return ux.GetByKey(s.Ref)
	}
	if s.UserID > 0 {
		return ux.GetByKey(s.UserID)
	}
	if s.U != nil && s.U.ID > 0 {
		return ux.GetByKey(s.U.ID)
	}
	return s
}
