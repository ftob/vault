package vault

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	NotDirty = iota
	Dirty
)

//
type Vault interface {
	Len() int
	Put(key string, value interface{})
	Get(key string) interface{}
	Keys() []string
}

//
type keys struct {
	mx sync.Mutex
	ks []string
}

//
func (k *keys) add(key string) {
	k.mx.Lock()
	defer k.mx.Unlock()

	k.ks = append(k.ks, key)
}

func (k *keys) touch(key string) {
	k.mx.Lock()
	defer k.mx.Unlock()

	var te string
	ks := make([]string, 0)

	for _, v := range k.ks {
		if v == key {
			te = v
		} else {
			ks = append(ks, v)
		}
	}

	k.ks = append(ks, te)
}

//
func (k *keys) shift() {
	k.mx.Lock()
	defer k.mx.Unlock()
	ks := make([]string, 0)
	for k, v := range k.ks {
		if k == 0 {
			continue
		}

		ks = append(ks, v)
	}
	k.ks = ks
}

//
type vault struct {
	store sync.Map
	keys  keys
	l     int
	cap   int
	dirty uint32
}

//
func (v *vault) Len() int {
	if v.isDirty() {
		v.l = v.len()
		v.unmarkDirty()
	}

	return v.l
}

//
func (v *vault) isDirty() bool {
	return v.dirty == Dirty
}

//
func (v *vault) markDirty() {
	atomic.StoreUint32(&v.dirty, Dirty)
}

///
func (v *vault) unmarkDirty() {
	atomic.StoreUint32(&v.dirty, NotDirty)
}

func (v vault) len() int {
	var i int
	v.store.Range(func(key, value interface{}) bool {
		i++
		return true
	})
	return i
}

func (v *vault) Put(key string, value interface{}) {
	v.store.Store(key, value)
	v.keys.add(key)
	v.markDirty()

	if v.Len() > v.cap {
		v.shiftStore()
	}
}

func (v *vault) shiftStore() {
	v.store.Delete(v.keys.ks[0])
	v.keys.shift()
	v.markDirty()
}

func (v *vault) Get(key string) interface{} {
	if e, ok := v.store.Load(key); ok {
		v.keys.touch(key)
		return e
	} else {
		return ""
	}
}

func (v *vault) Keys() []string {
	return v.keys.ks
}

func NewVault(cap int) Vault {
	v := &vault{}
	v.store, v.keys = sync.Map{}, newKeys(cap)
	v.cap, v.l = cap, 0
	v.markDirty()
	return v
}

func newKeys(cap int) keys {
	return keys{
		mx: sync.Mutex{},
		ks: make([]string, 0),
	}
}
