package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/crusttech/human/server/pkg/revisions"
	"github.com/crusttech/human/server/pkg/sql"
	"reflect"
	"time"
)

type (
	Attachment struct {
		ID          uint64         `json:"attachmentID,string"`
		TenantID    uint64         `json:"tenantID,string,omitempty"`
		ProjectID   uint64         `json:"projectID,string,omitempty"`
		NamespaceID uint64         `json:"namespaceID,string"`
		OwnerID     uint64         `json:"ownerID,string"`
		Kind        string         `json:"-"`
		Url         string         `json:"url,omitempty"`
		PreviewUrl  string         `json:"previewUrl,omitempty"`
		Name        string         `json:"name,omitempty"`
		Meta        AttachmentMeta `json:"meta"`
		CreatedAt   time.Time      `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
		DeletedAt   *time.Time     `json:"deletedAt,omitempty"`
	}

	AttachmentMeta struct {
		Original AttachmentFileMeta     `json:"original"`
		Preview  *AttachmentFileMeta    `json:"preview,omitempty"`
		Icon     *AttachmentIconMeta    `json:"icon,omitempty"`
		IconSvg  *AttachmentIconSvgMeta `json:"iconSvg,omitempty"`
	}

	AttachmentFileMeta struct {
		Size      int64                `json:"size"`
		Extension string               `json:"ext"`
		Mimetype  string               `json:"mimetype"`
		Image     *AttachmentImageMeta `json:"image,omitempty"`
	}

	AttachmentImageMeta struct {
		Width    int  `json:"width,omitempty"`
		Height   int  `json:"height,omitempty"`
		Animated bool `json:"animated"`
	}

	AttachmentIconMeta struct {
		Name    string `json:"name"`
		Library string `json:"library"`
	}

	AttachmentIconSvgMeta struct {
		Src string `json:"src"`
	}
)

func (r Attachment) Clone() *Attachment {
	dup := r
	dup.Meta = *r.Meta.Clone()

	if r.UpdatedAt != nil {
		v := *r.UpdatedAt
		dup.UpdatedAt = &v
	}

	if r.DeletedAt != nil {
		v := *r.DeletedAt
		dup.DeletedAt = &v
	}

	return &dup
}

func (r Attachment) Diff(cmp *Attachment) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &Attachment{}
	}
	if r.ID != cmp.ID {
		out = append(out, &revisions.Change{Key: "attachmentID", Old: []any{cmp.ID}, New: []any{r.ID}})
	}

	if r.TenantID != cmp.TenantID {
		out = append(out, &revisions.Change{Key: "tenantID", Old: []any{cmp.TenantID}, New: []any{r.TenantID}})
	}

	if r.ProjectID != cmp.ProjectID {
		out = append(out, &revisions.Change{Key: "projectID", Old: []any{cmp.ProjectID}, New: []any{r.ProjectID}})
	}

	if r.NamespaceID != cmp.NamespaceID {
		out = append(out, &revisions.Change{Key: "namespaceID", Old: []any{cmp.NamespaceID}, New: []any{r.NamespaceID}})
	}

	if r.OwnerID != cmp.OwnerID {
		out = append(out, &revisions.Change{Key: "ownerID", Old: []any{cmp.OwnerID}, New: []any{r.OwnerID}})
	}

	if r.Url != cmp.Url {
		out = append(out, &revisions.Change{Key: "url", Old: []any{cmp.Url}, New: []any{r.Url}})
	}

	if r.PreviewUrl != cmp.PreviewUrl {
		out = append(out, &revisions.Change{Key: "previewUrl", Old: []any{cmp.PreviewUrl}, New: []any{r.PreviewUrl}})
	}

	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	for _, c := range r.Meta.Diff(&cmp.Meta) {
		c.Key = "meta." + c.Key
		out = append(out, c)
	}

	if !reflect.DeepEqual(r.CreatedAt, cmp.CreatedAt) {
		out = append(out, &revisions.Change{Key: "createdAt", Old: []any{cmp.CreatedAt}, New: []any{r.CreatedAt}})
	}

	if !reflect.DeepEqual(r.UpdatedAt, cmp.UpdatedAt) {
		out = append(out, &revisions.Change{Key: "updatedAt", Old: []any{cmp.UpdatedAt}, New: []any{r.UpdatedAt}})
	}

	if !reflect.DeepEqual(r.DeletedAt, cmp.DeletedAt) {
		out = append(out, &revisions.Change{Key: "deletedAt", Old: []any{cmp.DeletedAt}, New: []any{r.DeletedAt}})
	}

	return out
}

func (r AttachmentMeta) Clone() *AttachmentMeta {
	dup := r
	dup.Original = *r.Original.Clone()

	if r.Preview != nil {
		dup.Preview = r.Preview.Clone()
	}

	if r.Icon != nil {
		dup.Icon = r.Icon.Clone()
	}

	if r.IconSvg != nil {
		dup.IconSvg = r.IconSvg.Clone()
	}

	return &dup
}

func (r AttachmentMeta) Diff(cmp *AttachmentMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AttachmentMeta{}
	}
	for _, c := range r.Original.Diff(&cmp.Original) {
		c.Key = "original." + c.Key
		out = append(out, c)
	}

	if (r.Preview == nil) != (cmp.Preview == nil) {
		out = append(out, &revisions.Change{Key: "preview", Old: []any{cmp.Preview}, New: []any{r.Preview}})
	} else if r.Preview != nil {
		for _, c := range r.Preview.Diff(cmp.Preview) {
			c.Key = "preview." + c.Key
			out = append(out, c)
		}
	}

	if (r.Icon == nil) != (cmp.Icon == nil) {
		out = append(out, &revisions.Change{Key: "icon", Old: []any{cmp.Icon}, New: []any{r.Icon}})
	} else if r.Icon != nil {
		for _, c := range r.Icon.Diff(cmp.Icon) {
			c.Key = "icon." + c.Key
			out = append(out, c)
		}
	}

	if (r.IconSvg == nil) != (cmp.IconSvg == nil) {
		out = append(out, &revisions.Change{Key: "iconSvg", Old: []any{cmp.IconSvg}, New: []any{r.IconSvg}})
	} else if r.IconSvg != nil {
		for _, c := range r.IconSvg.Diff(cmp.IconSvg) {
			c.Key = "iconSvg." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *AttachmentMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AttachmentMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AttachmentFileMeta) Clone() *AttachmentFileMeta {
	dup := r
	if r.Image != nil {
		dup.Image = r.Image.Clone()
	}

	return &dup
}

func (r AttachmentFileMeta) Diff(cmp *AttachmentFileMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AttachmentFileMeta{}
	}
	if r.Size != cmp.Size {
		out = append(out, &revisions.Change{Key: "size", Old: []any{cmp.Size}, New: []any{r.Size}})
	}

	if r.Extension != cmp.Extension {
		out = append(out, &revisions.Change{Key: "ext", Old: []any{cmp.Extension}, New: []any{r.Extension}})
	}

	if r.Mimetype != cmp.Mimetype {
		out = append(out, &revisions.Change{Key: "mimetype", Old: []any{cmp.Mimetype}, New: []any{r.Mimetype}})
	}

	if (r.Image == nil) != (cmp.Image == nil) {
		out = append(out, &revisions.Change{Key: "image", Old: []any{cmp.Image}, New: []any{r.Image}})
	} else if r.Image != nil {
		for _, c := range r.Image.Diff(cmp.Image) {
			c.Key = "image." + c.Key
			out = append(out, c)
		}
	}

	return out
}

func (r *AttachmentFileMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AttachmentFileMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AttachmentImageMeta) Clone() *AttachmentImageMeta {
	dup := r
	return &dup
}

func (r AttachmentImageMeta) Diff(cmp *AttachmentImageMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AttachmentImageMeta{}
	}
	if r.Width != cmp.Width {
		out = append(out, &revisions.Change{Key: "width", Old: []any{cmp.Width}, New: []any{r.Width}})
	}

	if r.Height != cmp.Height {
		out = append(out, &revisions.Change{Key: "height", Old: []any{cmp.Height}, New: []any{r.Height}})
	}

	if r.Animated != cmp.Animated {
		out = append(out, &revisions.Change{Key: "animated", Old: []any{cmp.Animated}, New: []any{r.Animated}})
	}

	return out
}

func (r *AttachmentImageMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AttachmentImageMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AttachmentIconMeta) Clone() *AttachmentIconMeta {
	dup := r
	return &dup
}

func (r AttachmentIconMeta) Diff(cmp *AttachmentIconMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AttachmentIconMeta{}
	}
	if r.Name != cmp.Name {
		out = append(out, &revisions.Change{Key: "name", Old: []any{cmp.Name}, New: []any{r.Name}})
	}

	if r.Library != cmp.Library {
		out = append(out, &revisions.Change{Key: "library", Old: []any{cmp.Library}, New: []any{r.Library}})
	}

	return out
}

func (r *AttachmentIconMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AttachmentIconMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func (r AttachmentIconSvgMeta) Clone() *AttachmentIconSvgMeta {
	dup := r
	return &dup
}

func (r AttachmentIconSvgMeta) Diff(cmp *AttachmentIconSvgMeta) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &AttachmentIconSvgMeta{}
	}
	if r.Src != cmp.Src {
		out = append(out, &revisions.Change{Key: "src", Old: []any{cmp.Src}, New: []any{r.Src}})
	}

	return out
}

func (r *AttachmentIconSvgMeta) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r AttachmentIconSvgMeta) Value() (driver.Value, error) { return json.Marshal(r) }

func ParseAttachmentMeta(ss []string) (p AttachmentMeta, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}
