package types

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type (
	// Content is a stored row in the `content` table
	Content struct {
		ID       uint64 `json:"id,string" db:"id"`
		ModuleID uint64 `json:"moduleID,string" db:"module_id"`
		Page     *Page  `json:"page,omitempty"`

		Fields types.JSONText `json:"fields,omitempty" db:"-"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	// ContentColumn is a stored row in the `content_column` table
	ContentColumn struct {
		ContentID uint64 `json:"contentID,string" db:"content_id"`
		Name      string `json:"name" db:"column_name"`
		Value     string `json:"value" db:"column_value"`
	}

	// Field - CRM input field definitions
	Field struct {
		Name     string `json:"name" db:"field_name"`
		Type     string `json:"type" db:"field_type"`
		Template string `json:"template,omitempty" db:"field_template"`
	}

	// Modules - CRM module definitions
	Module struct {
		ID     uint64         `json:"id,string" db:"id"`
		Name   string         `json:"name" db:"name"`
		Fields types.JSONText `json:"fields" db:"json"`

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
		Private   bool   `json:"isPrivate" db:"is_private"`
	}

	// Page - page structure
	Page struct {
		ID     uint64 `json:"id,string" db:"id"`
		SelfID uint64 `json:"selfID,string" db:"self_id"`

		ModuleID uint64  `json:"moduleID,string" db:"module_id"`
		Module   *Module `json:"module,omitempty" db:"-"`

		Title       string `json:"title" db:"title"`
		Description string `json:"description" db:"description"`

		Blocks types.JSONText `json:"blocks" db:"blocks"`

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
