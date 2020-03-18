package lru

import (
	"reflect"
	"sync"
	"time"

	"github.com/RocksonZeta/wrap/utils/osutil"
	"github.com/RocksonZeta/wrap/wraplog"
)

var pkg = reflect.TypeOf(Options{}).PkgPath()
var log = wraplog.Logger.Fork(pkg, "Lru")

type Lru struct {
	keys    map[string]*LruNode
	Head    *LruNode //头部是最老数据
	Tail    *LruNode //尾部是最新数据
	options Options
	lock    sync.Mutex
}

type Options struct {
	Ttl             int //seconds
	MaxAge          int //seconds
	MaxLength       int
	CleanupInterval int //seconds
}

//newLocalSessionUids ttl :seconds
func New(options Options) *Lru {
	log.Trace().Func("NewLru").Interface("options", options).Send()
	if options.CleanupInterval <= 0 {
		options.CleanupInterval = 30
	}
	r := &Lru{
		keys:    make(map[string]*LruNode),
		options: options,
	}
	osutil.Go(func() {
		for {
			<-time.NewTicker(time.Duration(options.CleanupInterval) * time.Second).C
			r.cleanup()
		}
	})
	return r
}

func (l *Lru) Add(key string, value interface{}) {
	l.lock.Lock()
	now := time.Now().Unix()
	var node *LruNode
	old, ok := l.keys[key]
	if ok {
		node = old
		node.Key = key
		node.Value = value
		node.detach()
	} else {
		node = &LruNode{Value: value, Ct: now, Ut: now}
	}
	l.keys[key] = node
	if l.Head == nil {
		l.Head = node
		l.Tail = node
	} else {
		l.Tail.append(node)
		l.Tail = node
	}
	if l.options.MaxLength > 0 && len(l.keys) > l.options.MaxLength {
		h := l.Head
		l.Head = h.Next
		h.detach()
		delete(l.keys, h.Key)
	}
	l.lock.Unlock()
}

func (l *Lru) Get(key string) interface{} {
	now := time.Now().Unix()
	l.lock.Lock()
	node, ok := l.keys[key]
	if !ok {
		l.lock.Unlock()
		return nil
	}
	if l.options.MaxAge > 0 && now-node.Ct > int64(l.options.MaxAge) || l.options.Ttl > 0 && now-node.Ut > int64(l.options.Ttl) {
		node.detach()
		delete(l.keys, key)
		l.lock.Unlock()
		return nil
	}
	node.Ut = time.Now().Unix()
	node.detach()
	l.Tail.append(node)
	l.Tail = node
	l.lock.Unlock()
	return node.Value
}

func (l *Lru) GetString(key string) string {
	v := l.Get(key)
	if nil == v {
		return ""
	}
	return v.(string)
}
func (l *Lru) GetInt(key string) int {
	v := l.Get(key)
	if nil == v {
		return 0
	}
	return v.(int)
}
func (l *Lru) Delete(key string) {
	l.lock.Lock()
	l.delete(key)
	l.lock.Unlock()
}
func (l *Lru) delete(key string) {
	node, ok := l.keys[key]
	if ok {
		node.detach()
		delete(l.keys, key)
	}
}

func (l *Lru) cleanup() {
	l.lock.Lock()
	cur := l.Head
	now := time.Now().Unix()
	//清理ttl
	if l.options.Ttl > 0 {
		for {
			if cur != nil && now-cur.Ut > int64(l.options.Ttl) {
				delete(l.keys, l.Head.Key)
				cur = l.Head.Next
				l.Head.detach()
				l.Head = cur
			} else {
				break
			}
		}
	}
	l.lock.Unlock()
}

type LruNode struct {
	Pre   *LruNode
	Next  *LruNode
	Key   string
	Value interface{}
	Ct    int64
	Ut    int64
}

func (n *LruNode) detach() {
	if n == nil {
		return
	}
	if n.Next != nil {
		n.Next.Pre = n.Pre
	}
	if n.Pre != nil {
		n.Pre.Next = n.Next
	}
	n.Pre = nil
	n.Next = nil
}

func (n *LruNode) append(LruNode *LruNode) {
	if LruNode == nil || n == nil {
		return
	}
	LruNode.Next = n
	n.Pre = LruNode
}
