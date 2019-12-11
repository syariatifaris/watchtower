package watchtower

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTowerWatch(t *testing.T) {
	t.Run("", func(t *testing.T) {
		got := New()
		assert.NotNil(t, got, "NewTowerWatch : assert not nil")
		assert.IsType(t, &watcher{}, got, "NewTowerWatch : assert is type")
		assert.Implements(t, (*WatchTower)(nil), got, "NewTowerWatch : assert implements")
	})
}

func TestWatcherAddWatchObject(t *testing.T) {
	type args struct {
		fixables []Fixable
	}
	tests := []struct {
		name  string
		model *watcher
		args  args
	}{
		{
			model: &watcher{},
			args: args{
				fixables: []Fixable{
					Fixable{
						Name:    "foo",
						Err:     "foobar",
						Healthy: func() bool { return true },
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.AddWatchObject(tt.args.fixables...)

			assert.Empty(t, tt.model.brokens, "AddWatchObject : assert empty")
		})
	}
}

func TestWatcherIsBadInfrastructure(t *testing.T) {
	tests := []struct {
		name  string
		model *watcher
		want  bool
	}{
		{
			name:  "case empty",
			model: &watcher{fixingInProcess: new(AtomicBool)},
			want:  false,
		},
		{
			name: "case not empty, but not required",
			model: &watcher{
				errmsgs: map[string]string{
					"foo": "bar",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.IsBadInfrastructure()

			assert.Equal(t, tt.want, got, "IsBadInfrastructure : assert equal")
		})
	}
}

func TestWatcherGetErrMessages(t *testing.T) {
	tests := []struct {
		name  string
		model *watcher
		want  []string
	}{
		{
			name:  "case empty",
			model: &watcher{},
			want:  []string{},
		},
		{
			name: "case not empty, but not required",
			model: &watcher{
				errmsgs: map[string]string{
					"foo": "bar",
				},
			},
			want: []string{"bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.GetErrMessages()

			assert.Equal(t, tt.want, got, "GetErrMessages : assert equal")
		})
	}
}
