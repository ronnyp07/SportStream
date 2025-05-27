package cclogger

import (
	"github.com/sts-solutions/base-code/cclog"
	"go.uber.org/zap/zapcore"
)

func toZapCoreFields(lfs ...LogField) []zapcore.Field {
	fields := make([]zapcore.Field, 0, len(lfs))

	for _, lf := range lfs {
		fields = append(fields, cclog.NewField(lf.Key, lf.Value))
	}

	return fields
}
