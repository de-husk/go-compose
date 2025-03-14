package compose_test

import (
	"math"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/de-husk/go-compose"
)

func TestChain_Math(t *testing.T) {
	val := compose.New(math.Floor, math.Sqrt, math.Abs).Compose(-1234)

	if val != math.Floor(math.Sqrt(math.Abs(-1234))) {
		t.Errorf("val should be %v", val)
	}
}

func TestMerge(t *testing.T) {
	var i atomic.Uint32

	c1 := compose.New(i.Add)

	c2 := compose.New(i.Add)

	c3 := c1.Merge(c2)

	v1 := c1.Compose(5)
	if v1 != 5 {
		t.Errorf("v1 should be 5 but was: %v", v1)
	}

	v2 := c1.Compose(5)
	if v2 != 10 {
		t.Errorf("v1 should be 10 but was: %v", v2)
	}

	v3 := c3.Compose(10)

	if v3 != 40 {
		t.Errorf("v3 should be 40 but was: %v", v3)
		return
	}
}

// Http handler chain tests
type TestResWriter struct{}

func (rw *TestResWriter) Header() http.Header {
	return http.Header{}
}

func (rw *TestResWriter) Write([]byte) (int, error) {
	return -1, nil
}

func (rw *TestResWriter) WriteHeader(statusCode int) {}

func TestChain_OneHandler(t *testing.T) {
	count := 0
	hf := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count++

			next.ServeHTTP(w, r)
		})
	}

	c := compose.New(hf)

	h := c.Compose(http.DefaultServeMux)

	if h == nil {
		t.Errorf("nil handler")
		return
	}

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// Call composed chain
	h.ServeHTTP(&TestResWriter{}, req)

	// Ensure handler func was called:
	if count != 1 {
		t.Errorf("expected count to be %v - but was %v", 1, count)
	}
}

func TestChain_TwoHandlers(t *testing.T) {
	count := 0
	oneSet := false
	twoSet := false

	one := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if oneSet {
				t.Errorf("one has already been called")
				return
			}

			if twoSet {
				t.Errorf("two has been called before one")
				return
			}

			count++
			oneSet = true

			next.ServeHTTP(w, r)
		})
	}

	two := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !oneSet {
				t.Errorf("two has been called before one")
				return
			}

			if twoSet {
				t.Errorf("two has already been called")
				return
			}

			count++
			twoSet = true

			next.ServeHTTP(w, r)
		})
	}

	c := compose.New(one, two)

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	h := c.Compose(http.DefaultServeMux)

	if h == nil {
		t.Errorf("nil handler")
		return
	}

	// Call composed chain
	h.ServeHTTP(&TestResWriter{}, req)

	// Ensure handler func was called:
	if count != 2 {
		t.Errorf("expected count to be %v - but was %v", 2, count)
	}

	if !oneSet {
		t.Errorf("one was not called")
	}

	if !twoSet {
		t.Errorf("two was not called")
	}
}

func TestNext(t *testing.T) {
	count := 0
	oneSet := false
	twoSet := false

	one := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if oneSet {
				t.Errorf("one has already been called")
				return
			}

			if twoSet {
				t.Errorf("two has been called before one")
				return
			}

			count++
			oneSet = true

			next.ServeHTTP(w, r)
		})
	}

	two := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !oneSet {
				t.Errorf("two has been called before one")
				return
			}

			if twoSet {
				t.Errorf("two has already been called")
				return
			}

			count++
			twoSet = true

			next.ServeHTTP(w, r)
		})
	}

	c := compose.New(one).Next(two)

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	h := c.Compose(http.DefaultServeMux)

	if h == nil {
		t.Errorf("nil handler")
		return
	}

	// Call composed chain
	h.ServeHTTP(&TestResWriter{}, req)

	// Ensure handler func was called:
	if count != 2 {
		t.Errorf("expected count to be %v - but was %v", 2, count)
	}

	if !oneSet {
		t.Errorf("one was not called")
	}

	if !twoSet {
		t.Errorf("two was not called")
	}
}

func one(h http.Handler) http.Handler {
	return h
}

func two(h http.Handler) http.Handler {
	return h
}

func TestAliasTypeInference(t *testing.T) {
	var hfs = []func(http.Handler) http.Handler{
		one,
		two,
	}
	compose.New(one, two)
	compose.New(hfs[0], hfs[1])
	compose.New(hfs...)
}
