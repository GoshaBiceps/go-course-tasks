// Задание: Корреляция логов и трейсов
//
// Добавь trace_id и span_id в slog-логи, если они присутствуют в context.
// По логу должно быть возможно найти соответствующий трейс.
//
// Ожидаемый результат (JSON-лог):
//   {"level":"INFO","msg":"handling request","trace_id":"abc123...","span_id":"def456...","user_id":"u-1"}
//   {"level":"ERROR","msg":"operation failed","trace_id":"abc123...","span_id":"def456...","error":"..."}

package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// traceFields извлекает trace_id и span_id из context (если есть активный span OTel).
//
// TODO: реализуй функцию traceFields(ctx context.Context) []any
// Шаги:
//   1. Получи SpanContext: span := trace.SpanFromContext(ctx)
//   2. Если span.SpanContext().IsValid() — добавь поля trace_id и span_id
//   3. Верни срез []any{"trace_id", "...", "span_id", "..."}
//   4. Если span невалидный — верни nil
//
// Подсказка: requires go.opentelemetry.io/otel/trace
// span.SpanContext().TraceID().String()
// span.SpanContext().SpanID().String()

func traceFields(ctx context.Context) []any {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return []any{
			"trace_id", span.SpanContext().TraceID().String(),
			"span_id", span.SpanContext().SpanID().String(),
		}
	}
	return nil
}

// logWithTrace логирует сообщение, автоматически добавляя trace_id/span_id из context.
//
// TODO: реализуй logWithTrace(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, args ...any)
// Объедини args и traceFields(ctx), затем вызови logger.Log(ctx, level, msg, allArgs...)

func logWithTrace(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, args ...any) {
	traceArgs := traceFields(ctx) // получаем наши значния из спанов

	allArgs := append(args, traceArgs...) // добавляем значнеи из нашего слайса в слайс аргументов

	logger.Log(ctx, level, msg, allArgs...) // записываем это все в логи

}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()

	// TODO: инициализируй OTel TracerProvider (аналогично task04)
	exporter, err := stdouttrace.New( // экспоретр
		stdouttrace.WithPrettyPrint(), // форматирование красиво
	)

	if err != nil { // если экспортер не создался
		log.Fatal(err) // печать ошибки и завершенеи программы
	}

	tp := sdktrace.NewTracerProvider( // все спаны котоыре будут создаваться через этот провайдер опбрабатывай батчером
		sdktrace.WithBatcher(exporter),
	)

	defer tp.Shutdown(ctx)

	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("token-service")

	// Пример использования (раскомментируй после реализации):
	ctx, span := tracer.Start(ctx, "handle-issue-token")
	defer span.End()

	logWithTrace(ctx, logger, slog.LevelInfo, "handling request", "user_id", "u-1")
	logWithTrace(ctx, logger, slog.LevelError, "operation failed", "error", "db timeout")

}
