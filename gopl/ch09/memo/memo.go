package memo

// Func is the type of the function to memoize.
type Func func(key string, done <-chan struct{}) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

type request struct {
	key      string
	response chan<- result
	done     <-chan struct{} // cancellation channel
}

// Memo implements simple memoization of a function.
type Memo struct {
	requests, cancels chan request
}

// New initialized memoizarion structure.
func New(f Func) *Memo {
	memo := &Memo{
		requests: make(chan request),
		cancels:  make(chan request),
	}
	go memo.server(f)
	return memo
}

// Get calculates new result for Func or returns one
// from the cache.
func (m *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	req := request{
		key:      key,
		response: response,
		done:     done,
	}
	m.requests <- req
	res := <-response
	select {
	case <-done:
		m.cancels <- req
	default:
	}
	return res.value, res.err
}

// Close closes internal channel for requests.
func (m *Memo) Close() {
	close(m.requests)
	close(m.cancels)
}

func (m *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range m.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key, req.done)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string, done <-chan struct{}) {
	e.res.value, e.res.err = f(key, done)
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	<-e.ready
	response <- e.res
}
