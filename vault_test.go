package vault

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewVault(t *testing.T) {
	type args struct {
		cap int
	}
	tests := []struct {
		name string
		args args
		want Vault
	}{
		struct {
			name string
			args args
			want Vault
		}{name: "positive", args: args{cap: 1}, want: &vault{
			store: sync.Map{},
			keys: keys{
				mx: sync.Mutex{},
				ks: make([]string, 1),
			},
			l:     0,
			cap:   1,
			dirty: 0,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVault(tt.args.cap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keys_add(t *testing.T) {
	type fields struct {
		mx sync.Mutex
		ks []string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		struct {
			name   string
			fields fields
			args   args
		}{name: "", fields: fields{
			mx: sync.Mutex{},
			ks: make([]string, 1),
		}, args: args{key: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := keys{
				mx: tt.fields.mx,
				ks: tt.fields.ks,
			}
			t.Log(k)
		})
	}
}

func Test_keys_deleteFirst(t *testing.T) {
	type fields struct {
		mx sync.Mutex
		ks []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		struct {
			name   string
			fields fields
		}{name: "", fields: fields{
			mx: sync.Mutex{},
			ks: make([]string, 1),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := keys{
				mx: tt.fields.mx,
				ks: tt.fields.ks,
			}
			t.Log(k)
		})
	}
}

func Test_vault_Get(t *testing.T) {
	type fields struct {
		store sync.Map
		keys  keys
		l     int
		cap   int
		dirty uint32
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		struct {
			name   string
			fields fields
			args   args
			want   interface{}
		}{name: "", fields: fields{
			store: sync.Map{},
			keys:  keys{},
			l:     0,
			cap:   1,
			dirty: 1,
		}, args: args{key: ""}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vault{
				store: tt.fields.store,
				keys:  tt.fields.keys,
				l:     tt.fields.l,
				cap:   tt.fields.cap,
				dirty: tt.fields.dirty,
			}
			if got := v.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vault_Keys(t *testing.T) {
	type fields struct {
		store sync.Map
		keys  keys
		l     int
		cap   int
		dirty uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		struct {
			name   string
			fields fields
			want   []string
		}{name: "", fields: fields{
			store: sync.Map{},
			keys: keys{
				mx: sync.Mutex{},
				ks: []string{},
			},
			l:     0,
			cap:   0,
			dirty: 0,
		}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vault{
				store: tt.fields.store,
				keys:  tt.fields.keys,
				l:     tt.fields.l,
				cap:   tt.fields.cap,
				dirty: tt.fields.dirty,
			}
			if got := v.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vault_Len(t *testing.T) {
	type fields struct {
		store sync.Map
		keys  keys
		l     int
		cap   int
		dirty uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		struct {
			name   string
			fields fields
			want   int
		}{name: "", fields: fields{
			store: sync.Map{},
			keys:  keys{},
			l:     0,
			cap:   0,
			dirty: 0,
		}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vault{
				store: tt.fields.store,
				keys:  tt.fields.keys,
				l:     tt.fields.l,
				cap:   tt.fields.cap,
				dirty: tt.fields.dirty,
			}
			if got := v.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vault_Put(t *testing.T) {
	type fields struct {
		store sync.Map
		keys  keys
		l     int
		cap   int
		dirty uint32
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		struct {
			name   string
			fields fields
			args   args
		}{name: "", fields: fields{
			store: sync.Map{},
			keys:  keys{},
			l:     0,
			cap:   0,
			dirty: 0,
		}, args: args{
			key:   "",
			value: nil,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vault{
				store: tt.fields.store,
				keys:  tt.fields.keys,
				l:     tt.fields.l,
				cap:   tt.fields.cap,
				dirty: tt.fields.dirty,
			}
			t.Log(v)
		})
	}
}

func Test_vault_isDirty(t *testing.T) {
	type fields struct {
		store sync.Map
		keys  keys
		l     int
		cap   int
		dirty uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		struct {
			name   string
			fields fields
			want   bool
		}{name: "", fields: fields{
			store: sync.Map{},
			keys:  keys{},
			l:     0,
			cap:   0,
			dirty: 1,
		}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vault{
				store: tt.fields.store,
				keys:  tt.fields.keys,
				l:     tt.fields.l,
				cap:   tt.fields.cap,
				dirty: tt.fields.dirty,
			}
			if got := v.isDirty(); got != tt.want {
				t.Errorf("isDirty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vault_len(t *testing.T) {
	type fields struct {
		store sync.Map
		keys  keys
		l     int
		cap   int
		dirty uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		struct {
			name   string
			fields fields
			want   int
		}{name: "", fields: fields{
			store: sync.Map{},
			keys:  keys{},
			l:     0,
			cap:   0,
			dirty: 0,
		}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vault{
				store: tt.fields.store,
				keys:  tt.fields.keys,
				l:     tt.fields.l,
				cap:   tt.fields.cap,
				dirty: tt.fields.dirty,
			}
			if got := v.len(); got != tt.want {
				t.Errorf("len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vault_markDirty(t *testing.T) {
	type fields struct {
		store sync.Map
		keys  keys
		l     int
		cap   int
		dirty uint32
	}
	tests := []struct {
		name   string
		fields fields
	}{
		struct {
			name   string
			fields fields
		}{name: "", fields: fields{
			store: sync.Map{},
			keys:  keys{},
			l:     0,
			cap:   0,
			dirty: 0,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vault{
				store: tt.fields.store,
				keys:  tt.fields.keys,
				l:     tt.fields.l,
				cap:   tt.fields.cap,
				dirty: tt.fields.dirty,
			}
			t.Log(v)
		})
	}
}

func Test_vault_unmarkDirty(t *testing.T) {
	type fields struct {
		store sync.Map
		keys  keys
		l     int
		cap   int
		dirty uint32
	}
	tests := []struct {
		name   string
		fields fields
	}{
		struct {
			name   string
			fields fields
		}{name: "", fields: fields{
			store: sync.Map{},
			keys:  keys{},
			l:     0,
			cap:   0,
			dirty: 0,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vault{
				store: tt.fields.store,
				keys:  tt.fields.keys,
				l:     tt.fields.l,
				cap:   tt.fields.cap,
				dirty: tt.fields.dirty,
			}
			t.Log(v)
		})
	}
}
