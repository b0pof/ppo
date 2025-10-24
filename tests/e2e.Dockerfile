FROM golang:1.24-alpine AS build

WORKDIR /project
COPY . .
RUN go mod download
RUN go build -o server ./cmd/service/main.go

FROM golang:1.24-alpine

WORKDIR /project
COPY --from=build /project /project

RUN mkdir -p /allure-results
ENV ALLURE_OUTPUT_PATH=/allure-results

EXPOSE 8080

CMD sh -c "\
    ./server & \
    SERVER_PID=$! && \
    sleep 5 && \
    go test -tags=e2e ./tests/e2e/... -v && \
    TEST_EXIT_CODE=$? && \
    kill $SERVER_PID && \
    exit $TEST_EXIT_CODE"
