package logger

import (
	"context"

	"kudago/internal/http/utils"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Logger struct {
	Logger *zap.SugaredLogger
}

func NewLogger() (*Logger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{Logger: zapLogger.Sugar()}, nil
}

func (l *Logger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	requestID := utils.GetRequestIDFromContext(ctx)
	l.Logger.Infow("Query",
		"request_id", requestID,
		"sql", data.SQL,
		"args", data.Args,
	)
	return ctx
}

func (l *Logger) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	requestID := utils.GetRequestIDFromContext(ctx)
	if data.Err != nil {
		l.Logger.Errorw("Query failed",
			"request_id", requestID,
			"commandTag", data.CommandTag,
			"args", data.Err,
		)
	}
}

func (l *Logger) Error(err error, args ...interface{}) {
	l.Logger.Error("failed : %v, args: %v", zap.Error(err), args)
}
