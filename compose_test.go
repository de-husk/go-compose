package compose_test

import (
	"math"
	"net/http"
	"testing"

	"github.com/de-husk/go-compose"
)

func TestChain_Math(t *testing.T) {
	val := compose.NewChain(math.Floor, math.Sqrt, math.Abs).Compose(-1234)

	if val != math.Floor(math.Sqrt(math.Abs(-1234))) {
		t.Errorf("val should be %v", val)
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

	c := compose.NewChain(hf)

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

	c := compose.NewChain(one, two)

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

	c := compose.NewChain(one)

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	c.Next(two)

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
