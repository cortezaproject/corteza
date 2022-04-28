package crs

import (
	"context"
	"fmt"

	ccrs "github.com/cortezaproject/corteza-server/compose/crs"
	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type (
	Connector struct {
		DB connection

		// Logger for connection
		Logger *zap.Logger

		Dialect goqu.DialectWrapper
		// ---

		modelHandlers map[uint64]*crs
		capabilities  capabilities.Set
	}
)

func CRSConnector(db connection, logger *zap.Logger, di goqu.DialectWrapper, cc capabilities.Set) *Connector {
	return &Connector{
		DB:            db,
		Logger:        logger,
		Dialect:       di,
		modelHandlers: make(map[uint64]*crs),
		capabilities:  cc,
	}
}

func (c *Connector) Capabilities() capabilities.Set {
	return c.capabilities
}

func (c *Connector) Can(capabilities ...capabilities.Capability) bool {
	return c.Capabilities().IsSuperset(capabilities...)
}

func (c *Connector) Close(ctx context.Context) error {
	return nil
}

func (c *Connector) CreateRecords(ctx context.Context, m *data.Model, rr ...*types.Record) error {
	h, err := c.getModelHandlerE(m)
	if err != nil {
		return err
	}

	return h.Create(ctx, rr...)
}

func (c *Connector) SearchRecords(ctx context.Context, sch *data.Model, filter any) (ccrs.Iterator, error) {
	return nil, nil
}

// ---

func (c *Connector) RemoveModel(context.Context, *data.Model, ...*data.Model) error {
	return nil
}

func (c *Connector) AlterModel(ctx context.Context, old *data.Model, new *data.Model) error {
	return nil
}

func (c *Connector) AlterModelAttribute(ctx context.Context, sch *data.Model, old data.Attribute, new data.Attribute, trans ...ccrs.TransformationFunction) error {
	return nil
}

func (c *Connector) Models(ctx context.Context) (data.ModelSet, error) {
	// @todo...
	return nil, nil
}

func (c *Connector) AddModel(ctx context.Context, m *data.Model, mm ...*data.Model) error {
	for _, model := range append(mm, m) {
		if h := c.getModelHandler(model); h != nil {
			return fmt.Errorf("cannot add model: already exists")
		}

		_, err := c.addModelHandler(ctx, model)
		if err != nil {
			return err
		}
	}

	return nil
}

// ---

func (c *Connector) getModelHandler(m *data.Model) *crs {
	return c.modelHandlers[m.ResourceID]
}

func (c *Connector) getModelHandlerE(m *data.Model) (*crs, error) {
	h := c.modelHandlers[m.ResourceID]
	if h == nil {
		return nil, fmt.Errorf("cannot process model for resource %d: not supported", m.ResourceID)
	}
	return h, nil
}

func (c *Connector) addModelHandler(ctx context.Context, m *data.Model) (*crs, error) {
	h := CRS(m, c.DB)
	c.modelHandlers[m.ResourceID] = h

	return h, nil
}
