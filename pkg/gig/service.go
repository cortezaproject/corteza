package gig

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/id"
)

var (
	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}

	// wrapper around id.Next() that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}

	gigStore = make(map[uint64]Gig)
)

type (
	Service interface {
		Create(context.Context, UpdatePayload) (Gig, error)
		Read(context.Context, uint64) (Gig, error)
		Update(context.Context, Gig, UpdatePayload) (Gig, error)
		State(context.Context, Gig) (out interface{}, err error)
		Tasks(context.Context) (out TaskDefSet)
		Workers(context.Context) (out WorkerDefSet)
		Delete(context.Context, Gig) (Gig, error)
		Undelete(context.Context, Gig) (Gig, error)

		Prepare(context.Context, Gig) (Gig, error)
		Exec(context.Context, Gig) (Gig, error)
		Output(context.Context, Gig) (SourceWrapSet, Gig, error)
		OutputAll(context.Context, Gig) (SourceSet, Gig, error)
		OutputSpecific(context.Context, Gig, uint64) (Source, Gig, error)
		Cleanup(context.Context, Gig) (Gig, error)

		SetSources(context.Context, Gig, SourceWrapSet, ...Decoder) (Gig, error)
		AddSources(context.Context, Gig, SourceWrapSet, ...Decoder) (Gig, error)
		RemoveSources(context.Context, Gig, ...uint64) (Gig, error)
		SetDecoders(context.Context, Gig, ...Decoder) (Gig, error)
		SetPreprocessors(context.Context, Gig, ...Preprocessor) (Gig, error)
		SetPostprocessors(context.Context, Gig, ...Postprocessor) (Gig, error)
		Complete(context.Context, Gig) (Gig, error)
	}

	service struct {
		// @todo
		opt map[string]interface{}
	}
)

var (
	gSvc Service
)

// Initialize a new gig service
func NewService(opt map[string]interface{}) Service {
	if gSvc == nil {
		gSvc = &service{
			opt: opt,
		}
	}

	return gSvc
}

func (svc *service) Create(ctx context.Context, pl UpdatePayload) (g Gig, err error) {
	if pl.Worker == nil {
		err = fmt.Errorf("unable to create gig: worker not defined")
		return
	}

	g = newGig(pl.Worker)
	g.CompleteOn = pl.CompleteOn

	decoders := pl.Decode
	if len(pl.Sources) > 0 {
		g, err = setSources(ctx, g, pl.Sources, decoders...)
		if err != nil {
			return
		}
	}

	if len(pl.Preprocess) > 0 {
		g, err = setPreprocessors(ctx, g, pl.Preprocess)
		if err != nil {
			return
		}
	}

	if len(pl.Postprocess) > 0 {
		g, err = setPostprocessors(ctx, g, pl.Postprocess)
		if err != nil {
			return
		}
	}

	return updateGig(ctx, g)
}

func (svc *service) Update(ctx context.Context, old Gig, pl UpdatePayload) (g Gig, err error) {
	err = wrapValidationErr(old, func(old Gig) error {
		if pl.Worker != nil {
			return fmt.Errorf("worker can not be changed")
		}

		return baseChangeValidation(old)
	})
	if err != nil {
		return
	}

	g = old
	g.CompleteOn = pl.CompleteOn

	var implicit []Decoder
	if len(pl.Decode) > 0 {
		var explicit []Decoder
		for _, _d := range pl.Decode {
			d := _d

			if d.Source() != 0 {
				explicit = append(explicit, d)
			} else {
				implicit = append(implicit, d)
			}
		}

		g, err = setDecoders(ctx, g, explicit)
		if err != nil {
			return
		}
	}

	if len(pl.Sources) > 0 {
		g, err = setSources(ctx, g, pl.Sources, implicit...)
		if err != nil {
			return
		}
	}

	if len(pl.Preprocess) > 0 {
		g, err = setPreprocessors(ctx, g, pl.Preprocess)
		if err != nil {
			return
		}
	}

	if len(pl.Postprocess) > 0 {
		g, err = setPostprocessors(ctx, g, pl.Postprocess)
		if err != nil {
			return
		}
	}

	g.UpdatedAt = now()

	return updateGig(ctx, g)
}

func (svc *service) Read(ctx context.Context, id uint64) (Gig, error) {
	return getGig(ctx, id)
}

func (svc *service) AddSources(ctx context.Context, old Gig, sources SourceWrapSet, decoders ...Decoder) (g Gig, err error) {
	err = wrapValidationErr(old, baseChangeValidation)
	if err != nil {
		return
	}

	g, err = setSources(ctx, old, append(ToSourceWrap(old.Sources...), sources...), decoders...)
	if err != nil {
		return
	}

	return updateGig(ctx, g)
}

func (svc *service) SetSources(ctx context.Context, old Gig, sources SourceWrapSet, decoders ...Decoder) (g Gig, err error) {
	err = wrapValidationErr(old, baseChangeValidation)
	if err != nil {
		return
	}

	g, err = setSources(ctx, old, sources, decoders...)
	if err != nil {
		return
	}

	return updateGig(ctx, g)
}

func (svc *service) RemoveSources(ctx context.Context, old Gig, sources ...uint64) (g Gig, err error) {
	err = wrapValidationErr(old, baseChangeValidation)
	if err != nil {
		return
	}

	return removeSources(ctx, old, sources...)
}

func (svc *service) SetDecoders(ctx context.Context, old Gig, decoders ...Decoder) (g Gig, err error) {
	err = wrapValidationErr(old, baseChangeValidation)
	if err != nil {
		return
	}

	return setDecoders(ctx, old, decoders)
}

func (svc *service) SetPreprocessors(ctx context.Context, old Gig, preprocessors ...Preprocessor) (g Gig, err error) {
	err = wrapValidationErr(old, baseChangeValidation)
	if err != nil {
		return
	}

	return setPreprocessors(ctx, old, preprocessors)
}

func (svc *service) SetPostprocessors(ctx context.Context, old Gig, postprocessors ...Postprocessor) (g Gig, err error) {
	err = wrapValidationErr(old, baseChangeValidation)
	if err != nil {
		return
	}

	return setPostprocessors(ctx, old, postprocessors)
}

func (svc *service) Cleanup(ctx context.Context, old Gig) (g Gig, err error) {
	g, err = svc.cleanup(ctx, old)
	if err != nil {
		return
	}
	g.CompletedAt = now()
	return updateGig(ctx, g)
}

func (svc *service) cleanup(ctx context.Context, old Gig) (g Gig, err error) {
	g, err = removeSources(ctx, old)
	if err != nil {
		return
	}

	if g.Worker == nil {
		return
	}

	if err = g.Worker.Cleanup(ctx); err != nil {
		return
	}

	return
}

func (svc *service) Complete(ctx context.Context, old Gig) (g Gig, err error) {
	if old.CompletedAt != nil {
		err = fmt.Errorf("cannot complete gig %d: already completed", old.ID)
		return
	}

	g, err = svc.cleanup(ctx, old)
	if err != nil {
		return
	}

	g.CompletedAt = now()

	return updateGig(ctx, g)
}

func (svc *service) Prepare(ctx context.Context, old Gig) (g Gig, err error) {
	err = (func() error {
		if old.Worker == nil {
			return fmt.Errorf("worker is not defined")
		}
		if old.PreparedAt != nil {
			return fmt.Errorf("already prepared for execution")
		}
		return nil
	})()
	if err != nil {
		err = fmt.Errorf("cannot prepare gig %d: %w", old.ID, err)
		return
	}

	g = old
	g.PreparedAt = now()

	// Do decoding
	sources, err := runDecoders(ctx, old.Sources)
	if err != nil {
		return
	}

	// Do the worker preparations
	err = g.Worker.Prepare(ctx, sources...)
	if err != nil {
		return
	}

	err = g.Worker.Preprocess(ctx, old.Preprocess...)
	return
}

func (svc *service) Exec(ctx context.Context, old Gig) (g Gig, err error) {
	// Do the processing
	var (
		output SourceSet
		meta   WorkMeta
	)
	g = old

	// Prepare for exec in case it was skipped
	if g.PreparedAt == nil {
		g, err = svc.Prepare(ctx, old)
		if err != nil {
			return
		}
	}

	// Run
	g = gigExecStarted(g)
	output, meta, err = g.Worker.Exec(ctx)
	g, output, err = gigExecFinished(g, output, meta, err)
	if err != nil {
		return
	}

	g.Output = output

	// Do the postprocessing
	for _, pp := range g.Postprocess {
		output, g.Status.Meta, err = pp.Postprocess(ctx, g.Status.Meta, output)
		if err != nil {
			return
		}
	}

	g.Output = output

	if g.CompleteOn == OnExec {
		g, err = svc.cleanup(ctx, g)
		if err != nil {
			return
		}
		g.CompletedAt = now()
	}

	return updateGig(ctx, g)
}

func (svc *service) Output(ctx context.Context, old Gig) (out SourceWrapSet, g Gig, err error) {
	tmp := old.Output
	g = old

	if old.CompletedAt != nil {
		err = fmt.Errorf("unable to get output for gig %d: already completed", old.ID)
		return
	}

	if len(tmp) == 0 {
		return
	}

	out = ToSourceWrap(tmp...)
	return
}

func (svc *service) OutputAll(ctx context.Context, old Gig) (out SourceSet, g Gig, err error) {
	out = old.Output
	g = old

	if old.CompletedAt != nil {
		err = fmt.Errorf("unable to get output for gig %d: already completed", old.ID)
		return
	}

	if old.CompleteOn == OnOutput {
		g, err = svc.cleanup(ctx, old)
		if err != nil {
			return
		}
		g.CompletedAt = now()
		g, err = updateGig(ctx, g)
		if err != nil {
			return
		}
	}

	return
}

func (svc *service) OutputSpecific(ctx context.Context, old Gig, sourceID uint64) (out Source, g Gig, err error) {
	if old.CompletedAt != nil {
		err = fmt.Errorf("unable to get output for gig %d: already completed", old.ID)
		return
	}

	g = old
	tmp := g.Output.GetByID(sourceID)

	if tmp == nil {
		err = fmt.Errorf("unable to get the output for gig: %d: source not found: %d", g.ID, sourceID)
		return
	}

	out = *tmp
	return
}

func (svc *service) State(ctx context.Context, old Gig) (out interface{}, err error) {
	if old.CompletedAt != nil {
		err = fmt.Errorf("unable to get state for gig %d: already completed", old.ID)
		return
	}

	return old.Worker.State(ctx)
}

func (svc *service) Tasks(_ context.Context) (out TaskDefSet) {
	out = append(out, decoderDefinitions()...)
	out = append(out, preprocessorDefinitions()...)
	out = append(out, postprocessorDefinitions()...)

	return
}

func (svc *service) Workers(_ context.Context) (out WorkerDefSet) {
	return workerDefinitions()
}

func (svc *service) Delete(ctx context.Context, old Gig) (g Gig, err error) {
	old.DeletedAt = now()
	return updateGig(ctx, old)
}

func (svc *service) Undelete(ctx context.Context, old Gig) (g Gig, err error) {
	old.DeletedAt = nil
	return updateGig(ctx, old)
}

func baseChangeValidation(g Gig) error {
	if g.PreparedAt != nil {
		return fmt.Errorf("gig already prepared")
	}

	if g.CompletedAt != nil {
		return fmt.Errorf("gig already completed")
	}
	return nil
}

func wrapValidationErr(g Gig, v func(Gig) error) error {
	err := v(g)
	if err != nil {
		return fmt.Errorf("unable to update gig %d: %w", g.ID, err)
	}
	return nil
}

func setDecoders(ctx context.Context, old Gig, decoders []Decoder) (g Gig, err error) {
	g = old

	sourceMap := make(map[uint64][]Decoder)
	for _, _d := range decoders {
		d := _d

		src := g.Sources.GetByID(d.Source())
		if src == nil {
			err = fmt.Errorf("unknown source: %d", d.Source())
			return
		}

		sourceMap[d.Source()] = append(sourceMap[d.Source()], d)
	}

	for _, src := range g.Sources {
		if dd, ok := sourceMap[src.ID()]; ok {
			src.SetDecoders(dd...)
		}
	}

	return
}

func removeSources(ctx context.Context, old Gig, sources ...uint64) (g Gig, err error) {
	g = old

	var rm SourceSet
	var keep SourceSet

	if len(sources) > 0 {
		rmSet := make(map[uint64]bool)
		for _, sID := range sources {
			rmSet[sID] = true
		}

		for _, src := range g.Sources {
			if rmSet[src.ID()] {
				rm = append(rm, src)
			} else {
				keep = append(keep, src)
			}
		}
	} else {
		rm = g.Sources
	}

	err = cleanupSources(ctx, rm...)
	if err != nil {
		return
	}

	g.Sources = keep
	return
}

func setPreprocessors(ctx context.Context, old Gig, preprocessors []Preprocessor) (g Gig, err error) {
	g = old
	g.Preprocess = preprocessors

	return
}

func setPostprocessors(ctx context.Context, old Gig, postprocessors []Postprocessor) (g Gig, err error) {
	g = old
	g.Postprocess = postprocessors

	return
}

func setSources(ctx context.Context, old Gig, wraps SourceWrapSet, decoders ...Decoder) (g Gig, err error) {
	g = old

	crtSrcIndex := mapSources(g.Sources)
	sources := make(SourceSet, 0, len(g.Sources))
	cleanupSources := make(SourceSet, 0, len(g.Sources)/2)

	// Preserve unchanged sources
	for _, _src := range g.Sources {
		src := _src

		if wraps.Has(src.ID()) {
			sources = append(sources, src)
		} else {
			cleanupSources = append(cleanupSources, src)
		}
	}

	// New sources
	defaultSrc := make([]uint64, 0, 2)
	var src SourceSet
	for _, _w := range wraps {
		w := _w

		if _, ok := crtSrcIndex[w.ID]; !ok {
			src, err = FromSourceWrap(ctx, w)
			if err != nil {
				return
			}

			sources = append(sources, src...)
			for _, s := range src {
				defaultSrc = append(defaultSrc, s.ID())
			}
		}
	}

	// Cleanup old stuff
	for _, src := range cleanupSources {
		if err = src.Cleanup(); err != nil {
			return
		}
	}

	g.Sources = sources

	// Process decoders when provided
	if len(decoders) > 0 {
		explicit := make([]Decoder, 0, len(decoders))
		implicit := make([]Decoder, 0, len(decoders))

		for _, _d := range decoders {
			d := _d

			if d.Source() != 0 {
				explicit = append(explicit, d)
			} else {
				for _, src := range defaultSrc {
					implicit = append(implicit, d.Clone(src))
				}
			}
		}

		g, err = setDecoders(ctx, g, append(explicit, implicit...))
		if err != nil {
			return
		}
	}

	return
}

func cleanupSources(ctx context.Context, sources ...Source) (err error) {
	for _, s := range sources {
		if err = s.Cleanup(); err != nil {
			return
		}
	}

	return
}

func gigExecStarted(old Gig) Gig {
	old.Status.StartedAt = now()
	return old
}

func gigExecFinished(old Gig, output SourceSet, meta WorkMeta, err error) (Gig, SourceSet, error) {
	if err != nil {
		return old, output, err
	}

	// Time stuff
	old.Status.CompletedAt = now()
	old.Status.Elapsed = old.Status.CompletedAt.Sub(*old.Status.StartedAt)

	// Meta stuff
	old.Status.Error = err
	old.Status.Meta = meta

	return old, output, err
}
