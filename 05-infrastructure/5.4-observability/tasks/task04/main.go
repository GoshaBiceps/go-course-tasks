// Задание: Базовый трейсинг с OpenTelemetry
//
// Подключи OpenTelemetry SDK с stdout-экспортером (для dev).
// Создай span для HTTP-запроса и вложенный span для операции сервиса.
//
// Для запуска нужны зависимости:
//   go get go.opentelemetry.io/otel
//   go get go.opentelemetry.io/otel/sdk/trace
//   go get go.opentelemetry.io/otel/exporters/stdout/stdouttrace
//
// Ожидаемый результат (в stdout — JSON с данными трейса):
//   {"Name":"process-token","SpanContext":{...},...}
//   {"Name":"GET /api/v1/tokens","SpanContext":{...},...}

package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// TODO: реализуй initTracer() (*sdktrace.TracerProvider, error)
// Шаги:
//  1. Создай stdout exporter: stdouttrace.New(stdouttrace.WithPrettyPrint())
//  2. Создай TracerProvider с BatchSpanProcessor
//  3. Установи глобальный провайдер: otel.SetTracerProvider(tp)
//  4. Верни TracerProvider
func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New( // так создали экспортер отвечает куда отправить готовые данные о спанах
		stdouttrace.WithPrettyPrint(), // формотрирут джейсон красиво и читабельно
	)

	if err != nil { // проверяем если экспортер не создался
		return nil, err
	}

	tp := sdktrace.NewTracerProvider( // Все span'ы, которые будут создаваться через этот provider, обрабатывай батчером и отправляй через наш exporter».
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp) //установили глобально для нашего приложения

	return tp, nil
}

// TODO: вызови initTracer() в main
// Добавь defer tp.Shutdown(ctx) для сброса оставшихся span'ов

// TODO: создай tracer: otel.Tracer("token-service")

// TODO: в функции handleRequest(ctx context.Context, userID string):
//  1. Создай span "GET /api/v1/tokens": ctx, span := tracer.Start(ctx, "GET /api/v1/tokens")
//  2. Внутри вызови processToken(ctx, userID)
//  3. Закрой span: defer span.End()
func handleRequest(ctx context.Context, userID string) {
	tracer := otel.Tracer("token-service") //

	ctx, span := tracer.Start(ctx, "GET /api/v1/tokens")
	defer span.End()

	processToken(ctx, userID)

}

// TODO: в функции processToken(ctx context.Context, userID string):
//  1. Создай вложенный span "process-token"
//  2. Добавь атрибут: span.SetAttributes(attribute.String("user.id", userID))
//  3. Закрой span
func processToken(ctx context.Context, userID string) {
	tracer := otel.Tracer("token-service")

	ctx, childSpan := tracer.Start(ctx, "process-token")
	defer childSpan.End()

	childSpan.SetAttributes(attribute.String("user.id", userID))
}

func main() {
	ctx := context.Background()

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(ctx)

	handleRequest(
		ctx,
		"user-123",
	)

}
