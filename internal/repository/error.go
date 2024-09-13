package repository

import (
	"github.com/t34-dev/go-utils/pkg/sys"
	"github.com/t34-dev/go-utils/pkg/sys/codes"
)

var (
	ErrFailedCreateTable = sys.NewCommonError("failed to create table", codes.Internal)
)
