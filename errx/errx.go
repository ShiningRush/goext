package errx

import (
	"bytes"
	"fmt"
)

type BatchErrors struct {
	errs []error
}

func (be *BatchErrors) HasError() bool {
	return be.Len() > 0
}

func (be *BatchErrors) Len() int {
	return len(be.errs)
}

func (be *BatchErrors) Error() string {
	buf := new(bytes.Buffer)
	for idx, v := range be.errs {
		buf.WriteString(fmt.Sprintf("\nerr[%d]: %s", idx, v))
	}
	return buf.String()
}

func (be *BatchErrors) Append(err error) {
	be.errs = append(be.errs, err)
}
