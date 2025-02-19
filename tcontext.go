package tcontext

import (
	"context"
	"time"
)

type Context[T any] struct {
	context.Context
}

func _[T any]() {
	var _ context.Context = (*Context[T])(nil)
}

type tcontextKey struct{}

// FromContext returns the tcontext.Context from ctx, if ctx contains a value set
// with tcontextKey{}.
// This allows us to use tcontext.Context in code that accepts and transforms
// context.Context, and still be able to access the data stored in the
// tcontext.Context.
//
// If ctx is not a tcontext.Context, or a child context of a tcontext.Context,
// it will return a tcontext.Context with zero value of Data.
func FromContext[T any](ctx context.Context) (tctx Context[T], ok bool) {
	tctx, ok = ctx.(Context[T])
	if ok {
		return tctx, true
	}
	_, ok = ctx.Value(tcontextKey{}).(T)
	if ok {
		tctx = Context[T]{Context: ctx}
		return tctx, true
	}
	var t T
	tctx = WithData(ctx, t)
	return tctx, false
}

// WithData returns a new tcontext.Context that carries the data.
func WithData[T any](ctx context.Context, data T) Context[T] {
	return Context[T]{
		Context: context.WithValue(ctx, tcontextKey{}, data),
	}
}

// Data returns the data stored in the tcontext.Context.
func (tc Context[T]) Data() T {
	return tc.Value(tcontextKey{}).(T)
}

// Reimplementation of context.Context functions

// Analogous to context.WithCancel, but returns a tcontext.Context.
func WithCancel[T any](parent Context[T]) (ctx Context[T], cancel context.CancelFunc) {
	child, cancel := context.WithCancel(parent.Context)
	ctx = WithData(child, parent.Data())
	return
}

// Analogous to context.WithCancelCause, but returns a tcontext.Context.
func WithCancelCause[T any](parent Context[T]) (ctx Context[T], cancel context.CancelCauseFunc) {
	child, cancel := context.WithCancelCause(parent.Context)
	ctx = WithData(child, parent.Data())
	return
}

// Analogous to context.WithoutCancel, but returns a tcontext.Context.
func WithoutCancel[T any](parent Context[T]) Context[T] {
	return WithData(context.WithoutCancel(parent.Context), parent.Data())
}

// Analogous to context.WithoutCancelCause, but returns a tcontext.Context.
func WithDeadline[T any](parent Context[T], d time.Time) (Context[T], context.CancelFunc) {
	child, cancel := context.WithDeadline(parent.Context, d)
	return WithData(child, parent.Data()), cancel
}

// Analogous to context.WithDeadlineCause, but returns a tcontext.Context.
func WithDeadlineCause[T any](parent Context[T], d time.Time, cause error) (Context[T], context.CancelFunc) {
	child, cancel := context.WithDeadlineCause(parent.Context, d, cause)
	return WithData(child, parent.Data()), cancel
}

// Analogous to context.WithTimeout, but returns a tcontext.Context.
func WithTimeout[T any](parent Context[T], timeout time.Duration) (Context[T], context.CancelFunc) {
	child, cancel := context.WithTimeout(parent.Context, timeout)
	return WithData(child, parent.Data()), cancel
}

// Analogous to context.WithTimeoutCause, but returns a tcontext.Context.
func WithTimeoutCause[T any](parent Context[T], timeout time.Duration, cause error) (Context[T], context.CancelFunc) {
	child, cancel := context.WithTimeoutCause(parent.Context, timeout, cause)
	return WithData(child, parent.Data()), cancel
}
