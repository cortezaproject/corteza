package types

import (
	"time"

	"database/sql/driver"
	"encoding/json"

	"github.com/jmoiron/sqlx/types"

	systemTypes "github.com/crusttech/crust/system/types"
)

type (
	// Record is a stored row in the `record` table
	Record struct {
		ID       uint64 `json:"recordID,string" db:"id"`
		ModuleID uint64 `json:"moduleID,string" db:"module_id"`

		User   *systemTypes.User `json:"user,omitempty" db:"-"`
		UserID uint64            `json:"userID,string" db:"user_id"`

		Page *Page `json:"page,omitempty"`

		Fields types.JSONText `json:"fields,omitempty" db:"json"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	// RecordColumn is a stored row in the `record_column` table
	RecordColumn struct {
		RecordID uint64   `json:"-" db:"record_id"`
		Name     string   `json:"name" db:"column_name"`
		Value    string   `json:"value" db:"column_value"`
		Related  []string `json:"related" db:"-"`
	}

	Related struct {
		RecordID uint64 `json:"-" db:"record_id"`
		Name     string `json:"-" db:"column_name"`

		// RelatedRecordID isn't necessarily a record ID (multiple-select anything goes options)
		RelatedRecordID string `json:"-" db:"rel_record_id"`
	}

	// Modules - CRM module definitions
	Module struct {
		ID     uint64         `json:"moduleID,string" db:"id"`
		Name   string         `json:"name" db:"name"`
		Fields ModuleFieldSet `json:"fields" db:"json"`

		Page *Page `json:"page,omitempty"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	// Modules - CRM module definitions
	ModuleField struct {
		ModuleID uint64 `json:"moduleID,string" db:"module_id"`
		Place    int    `json:"-" db:"place"`

		Kind      string `json:"kind" db:"kind"`
		Name      string `json:"name" db:"name"`
		Label     string `json:"label" db:"label"`
		HelpText  string `json:"helpText,omitempty" db:"help_text"`
		Default   string `json:"defaultValue,omitempty" db:"default_value"`
		MaxLength int    `json:"maxLength" db:"max_length"`

		Options types.JSONText `json:"options" db:"json"`

		Private  bool `json:"isPrivate" db:"is_private"`
		Required bool `json:"isRequired" db:"is_required"`
		Visible  bool `json:"isVisible" db:"is_visible"`
	}

	ModuleFieldSet []ModuleField

	// Page - page structure
	Page struct {
		ID     uint64 `json:"pageID,string" db:"id"`
		SelfID uint64 `json:"selfID,string" db:"self_id"`

		ModuleID uint64  `json:"moduleID,string" db:"module_id"`
		Module   *Module `json:"module,omitempty" db:"-"`

		Title       string `json:"title" db:"title"`
		Description string `json:"description" db:"description"`

		Blocks types.JSONText `json:"blocks" db:"blocks"`

		Children PageSet `json:"children,omitempty" db:"-"`

		Visible bool `json:"visible" db:"visible"`
		Weight  int  `json:"-" db:"weight"`
	}

	// Block - value of Page.Blocks ([]Block)
	Block struct {
		Title       string         `json:"title"`
		Description string         `json:"description"`
		Options     types.JSONText `json:"options"`
		Kind        string         `json:"kind"`
		X           int            `json:"x"`
		Y           int            `json:"y"`
		Width       int            `json:"width"`
		Height      int            `json:"height"`
	}
)

func (f *ModuleFieldSet) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, f)
	}
	return nil
}

func (f ModuleFieldSet) Value() (driver.Value, error) {
	return json.Marshal(f)
}
