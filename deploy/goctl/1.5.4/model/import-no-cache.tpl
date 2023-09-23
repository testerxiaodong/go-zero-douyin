import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	{{if .time}}"time"{{end}}

	"go-zero-douyin/common/xconst"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)
