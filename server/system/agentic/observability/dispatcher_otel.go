package observability

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type otelDispatcher struct {
	tracer       trace.Tracer
	tp           *sdktrace.TracerProvider
	mu           sync.Mutex
	events       map[string][]AgentEvent      // events waiting to be attached to their span
	pending      map[string][]AgentSpan       // child spans whose parent hasn't arrived yet
	spanContexts map[string]trace.SpanContext // maps our span IDs to the OTel span context created for them
}

func NewOtelDispatcher(tp *sdktrace.TracerProvider) *otelDispatcher {
	return &otelDispatcher{
		tracer:       tp.Tracer("corteza.agent"),
		tp:           tp,
		events:       make(map[string][]AgentEvent),
		pending:      make(map[string][]AgentSpan),
		spanContexts: make(map[string]trace.SpanContext),
	}
}

func (d *otelDispatcher) ID() string { return "otel" }

func (d *otelDispatcher) OnSpan(span AgentSpan) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if span.ParentID == "" {
		d.createSpan(span, context.Background())
	} else if sc, ok := d.spanContexts[span.ParentID]; ok {
		d.createSpan(span, trace.ContextWithSpanContext(context.Background(), sc))
	} else {
		d.pending[span.ParentID] = append(d.pending[span.ParentID], span)
	}

	return nil
}

func (d *otelDispatcher) OnEvent(event AgentEvent) error {
	d.mu.Lock()
	d.events[event.SpanID] = append(d.events[event.SpanID], event)
	d.mu.Unlock()
	return nil
}

func (d *otelDispatcher) createSpan(span AgentSpan, parentCtx context.Context) {
	buffered := d.events[span.ID]
	delete(d.events, span.ID)

	_, otelSpan := d.tracer.Start(parentCtx, span.Name,
		trace.WithTimestamp(span.StartedAt),
		trace.WithSpanKind(trace.SpanKindInternal),
	)

	// save before end — children need this to set their parent
	d.spanContexts[span.ID] = otelSpan.SpanContext()

	attrs := []attribute.KeyValue{
		attribute.String("agent.id", span.AgentID),
		attribute.String("user.id", span.UserID),
		attribute.String("conversation.id", span.ConversationID),
	}
	for k, v := range span.Attributes {
		attrs = append(attrs, toOtelAttr(k, v))
	}
	otelSpan.SetAttributes(attrs...)

	for _, ev := range buffered {
		otelSpan.AddEvent(ev.Event, trace.WithTimestamp(ev.Timestamp))
	}

	if span.Status == StatusError {
		otelSpan.SetStatus(codes.Error, errStr(span.Error))
	} else {
		otelSpan.SetStatus(codes.Ok, "")
	}

	otelSpan.End(trace.WithTimestamp(span.EndedAt))

	d.flushPending(span.ID)
}

func (d *otelDispatcher) flushPending(parentID string) {
	children := d.pending[parentID]
	delete(d.pending, parentID)

	sc := d.spanContexts[parentID]
	ctx := trace.ContextWithSpanContext(context.Background(), sc)
	for _, child := range children {
		d.createSpan(child, ctx)
	}
}

func (d *otelDispatcher) Flush() error {
	return d.tp.ForceFlush(context.Background())
}

func (d *otelDispatcher) Shutdown() error {
	return d.tp.Shutdown(context.Background())
}

func toOtelAttr(key string, val any) attribute.KeyValue {
	switch v := val.(type) {
	case string:
		return attribute.String(key, v)
	case int:
		return attribute.Int64(key, int64(v))
	case int64:
		return attribute.Int64(key, v)
	case float64:
		return attribute.Float64(key, v)
	case bool:
		return attribute.Bool(key, v)
	default:
		return attribute.String(key, fmt.Sprintf("%v", v))
	}
}
