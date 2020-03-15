package vault

import (
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
func (k keys) add(key string) {
	k.mx.Lock()
	defer k.mx.Unlock()

	k.ks = append(k.ks, key)
}

//
func (k keys) shift() {
	k.mx.Lock()
	defer k.mx.Unlock()

	copy(k.ks, k.ks[1:])
}

//
type vault struct {
	store  sync.Map
	keys   keys
	l      int
	cap int
	dirty  uint32
}

//
func (v vault) Len() int {
	if v.isDirty() {
		v.l = v.len()
		v.unmarkDirty()
	}

	return v.l
}

//
func (v vault) isDirty() bool {
	return v.dirty == Dirty
}

//
func (v vault) markDirty() {
	atomic.StoreUint32(&v.dirty, Dirty)
}

///
func (v vault) unmarkDirty() {
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

func (v vault) Put(key string, value interface{}) {
	defer v.markDirty()
	v.store.Store(key, value)

	if v.Len() >= v.cap {
		v.store.Delete(v.keys.ks[0])
		v.keys.add(key)
		v.keys.shift()
	}
}

func (v vault) Get(key string) interface{} {
	if e, ok := v.store.Load(key); ok {
		return e
	} else {
		return ""
	}
}

func (v vault) Keys() []string {
	return v.keys.ks
}

func NewVault(cap int) Vault {
	v := &vault{}
	v.store = sync.Map{}

	v.keys = keys{
		mx: sync.Mutex{},
		ks: make([]string, cap),
	}
	v.cap = cap
	v.l = 0
	v.markDirty()
	return v
}
