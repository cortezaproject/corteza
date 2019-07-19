package types

import (
	"reflect"
	"testing"
)

func TestUnreadSet_Merge(t *testing.T) {
	type args struct {
		in UnreadSet
	}
	tests := []struct {
		name string
		uu   UnreadSet
		args args
		want UnreadSet
	}{
		{
			name: "simple",
			uu:   UnreadSet{&Unread{ChannelID: 1, UserID: 1, Count: 2}},
			args: args{in: UnreadSet{&Unread{ChannelID: 1, UserID: 1, ThreadCount: 3, ThreadTotal: 4}}},
			want: UnreadSet{&Unread{ChannelID: 1, UserID: 1, Count: 2, ThreadCount: 3, ThreadTotal: 4}},
		},
		{
			name: "empty base",
			uu:   UnreadSet{},
			args: args{in: UnreadSet{&Unread{ChannelID: 1, UserID: 1, ThreadCount: 3, ThreadTotal: 4}}},
			want: UnreadSet{&Unread{ChannelID: 1, UserID: 1, Count: 0, ThreadCount: 3, ThreadTotal: 4}},
		},
		{
			name: "emmpt input",
			uu:   UnreadSet{&Unread{ChannelID: 1, UserID: 1, Count: 2}},
			args: args{in: nil},
			want: UnreadSet{&Unread{ChannelID: 1, UserID: 1, Count: 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uu.Merge(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
