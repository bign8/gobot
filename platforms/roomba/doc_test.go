package roomba

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestPush(t *testing.T) {
	var err []error

	// Case 1: no pointer - no error
	gobottest.Assert(t, err, push(nil, nil))

	// Case 2: pointer - no error
	gobottest.Assert(t, err, push(&err, nil))

	// Case 3: no pointer - error
	sample := errors.New("Example 1")
	gobottest.Assert(t, []error{sample}, push(nil, sample))

	// Case 4: pointer - error
	push(&err, sample)
	gobottest.Assert(t, []error{sample}, err)
}
