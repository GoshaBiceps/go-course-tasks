package main

import (
	"context"
	"fmt"
	"time"
)

type UnaryHandler func(ctx context.Context, req any) (any, error) // описывает сам обработчик

type UnaryInterceptor func(ctx context.Context, method string, req any, handler UnaryHandler) (any, error) // описывает функцию , котора стоит вокруг  обработчика

type GreetRequest struct {
	Name string
}

type GreetResponse struct {
	Message string
}

type greetService struct{}

func (s *greetService) Greet(ctx context.Context, req *GreetRequest) (*GreetResponse, error) {
	time.Sleep(50 * time.Millisecond)
	return &GreetResponse{Message: "Hello, " + req.Name + "!"}, nil
}

// TODO: Implement LoggingInterceptor.
// It should return a UnaryInterceptor that:
// 1. Records start time
// 2. Calls the handler
// 3. Prints method name, duration, and error (if any)
// 4. Returns the handler's result

func LoggingInterceptor() UnaryInterceptor {
	return func(
		ctx context.Context,
		method string,
		req any,
		handler UnaryHandler,
	) (any, error) {

		start := time.Now() //

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			fmt.Printf(
				"method=%s duration=%s error=%v\n",
				method,
				duration,
				err,
			)
		} else {
			fmt.Printf(
				"method=%s duration=%s\n",
				method,
				duration,
			)
		}

		return resp, err
	}
}

func callWithInterceptor(
	ctx context.Context,
	method string,
	req any,
	handler UnaryHandler,
	interceptor UnaryInterceptor,
) (any, error) {
	return interceptor(ctx, method, req, handler)
}

func main() {
	svc := &greetService{}
	interceptor := LoggingInterceptor()
	ctx := context.Background()

	handler := func(ctx context.Context, req any) (any, error) {
		return svc.Greet(ctx, req.(*GreetRequest))
	}

	resp, err := callWithInterceptor(ctx, "/greet.GreetService/Greet", &GreetRequest{Name: "Alice"}, handler, interceptor)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Response: %+v\n", resp.(*GreetResponse))
}
