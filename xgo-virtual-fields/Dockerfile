FROM public.ecr.aws/d3j8x8q7/olympus-base-go:latest

WORKDIR /app
COPY . .

ENV GOFLAGS=-buildvcs=false
ENV GOTOOLCHAIN=local

RUN go mod download

RUN GOBIN=/usr/local/bin go install github.com/jstemmer/go-junit-report/v2@v2.1.0

RUN go build ./cl/... ./cmd/... ./parser/... ./printer/... && go test -count=1 -run "^$" ./cl/ >/dev/null

CMD ["/bin/bash"]
