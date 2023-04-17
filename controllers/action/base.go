package action

import "context"

type Action interface {
	Execute(context.Context) (noChange bool, err error)
}
