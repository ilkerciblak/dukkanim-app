package response

import "fmt"

type DomainError struct {
	Code   string
	Detail string
}

func (d  DomainError) Error() string {

	return fmt.Sprintf("%s: %s", d.Code, d.Detail)
}

func (d *DomainError) WithDetail(detail string) *DomainError {
	d.Detail = detail
	return d
}
