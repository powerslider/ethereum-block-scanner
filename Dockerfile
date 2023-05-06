FROM golang:1.20-alpine as builder

ENV PROJECT_NAME ethereum-block-scanner
ENV BASE_DIR /go/src/github.com/powerslider/${PROJECT_NAME}
WORKDIR ${BASE_DIR}

RUN apk --no-cache add git ca-certificates

COPY go.mod go.sum ${BASE_DIR}/

RUN go mod download -x

COPY cmd ${BASE_DIR}/cmd
COPY docs ${BASE_DIR}/docs
COPY pkg ${BASE_DIR}/pkg

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /dist/${PROJECT_NAME} ./cmd/${PROJECT_NAME}/main.go

FROM alpine

ENV PROJECT_NAME ethereum-block-scanner
ENV BASE_DIR /go/src/github.com/powerslider/${PROJECT_NAME}

COPY --from=builder /dist .

CMD ["/ethereum-block-scanner"]
