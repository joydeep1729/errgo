// +build go1.13

package errors

import (
	stderrors "errors"
	"fmt"
	"reflect"
	"testing"
)

func TestErrorChainCompat(t *testing.T) {
	err := stderrors.New("error that gets wrapped")
	wrapped := Wrap(err, "wrapped up")
	if !stderrors.Is(wrapped, err) {
		t.Errorf("Wrap does not support Go 1.13 error chains")
	}
}

func TestIs(t *testing.T) {
	err := New("test")

	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with stack",
			args: args{
				err:    WithStack(err),
				target: err,
			},
			want: true,
		},
		{
			name: "with message",
			args: args{
				err:    WithMessage(err, "test"),
				target: err,
			},
			want: true,
		},
		{
			name: "with message format",
			args: args{
				err:    WithMessagef(err, "%s", "test"),
				target: err,
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: err,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

type customErr struct {
	msg string
}

func (c customErr) Error() string { return c.msg }

func TestAs(t *testing.T) {
	var err = customErr{msg: "test message"}

	type args struct {
		err    error
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with stack",
			args: args{
				err:    WithStack(err),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "with message",
			args: args{
				err:    WithMessage(err, "test"),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "with message format",
			args: args{
				err:    WithMessagef(err, "%s", "test"),
				target: new(customErr),
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: new(customErr),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := As(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("As() = %v, want %v", got, tt.want)
			}

			ce := tt.args.target.(*customErr)
			if !reflect.DeepEqual(err, *ce) {
				t.Errorf("set target error failed, target error is %v", *ce)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	err := New("test")

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "with stack",
			args: args{err: WithStack(err)},
			want: err,
		},
		{
			name: "with message",
			args: args{err: WithMessage(err, "test")},
			want: err,
		},
		{
			name: "with message format",
			args: args{err: WithMessagef(err, "%s", "test")},
			want: err,
		},
		{
			name: "std errors compatibility",
			args: args{err: fmt.Errorf("wrap: %w", err)},
			want: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unwrap(tt.args.err); !reflect.DeepEqual(err, tt.want) {
				t.Errorf("Unwrap() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestUnwrapWithOuter(t *testing.T) {
	err := New("test")

	type args struct {
		err error
	}
	tests := []struct {
		name      string
		args      args
		wantInner error
		wantOuter error
	}{
		{
			name:      "with stack",
			args:      args{err: WithStack(err)},
			wantInner: err,
			wantOuter: WithStack(err),
		},
		{
			name:      "with message",
			args:      args{err: WithMessage(err, "test message")},
			wantInner: err,
			wantOuter: WithMessage(err, "test message"),
		},
		{
			name:      "with message format",
			args:      args{err: WithMessagef(err, "%s", "test fmt")},
			wantInner: err,
			wantOuter: WithMessagef(err, "%s", "test fmt"),
		},
		{
			name:      "std errors compatibility",
			args:      args{err: fmt.Errorf("wrap: %w", err)},
			wantInner: err,
			wantOuter: fmt.Errorf("wrap: %w", err),
		},
		{
			name:      "no inner error",
			args:      args{err: err}, // 'err' is just a fundamental error, has no Unwrap()
			wantInner: nil,
			wantOuter: err,
		},
		{
			name:      "nil error",
			args:      args{err: nil},
			wantInner: nil,
			wantOuter: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInner, gotOuter := UnwrapWithOuter(tt.args.err)
			
			// Compare inner
			if !reflect.DeepEqual(gotInner, tt.wantInner) {
				t.Errorf("UnwrapWithOuter() gotInner = %v, want %v", gotInner, tt.wantInner)
			}
			
			// Compare outer. For fmt.Errorf, DeepEqual might fail if we recreate it, so we compare Error strings if they aren't nil
			if gotOuter == nil || tt.wantOuter == nil {
				if gotOuter != tt.wantOuter {
					t.Errorf("UnwrapWithOuter() gotOuter = %v, want %v", gotOuter, tt.wantOuter)
				}
			} else if gotOuter.Error() != tt.wantOuter.Error() {
				t.Errorf("UnwrapWithOuter() gotOuter = %q, want %q", gotOuter.Error(), tt.wantOuter.Error())
			}
		})
	}
}
