FROM public.ecr.aws/d3j8x8q7/olympus-base-go:latest

WORKDIR /app
COPY . .

ENV GOFLAGS=-buildvcs=false

# TestLua runs scripts from the lua-tests submodule. Fetch it at the commit the
# repository pins when the build context doesn't include it.
RUN if [ ! -f lua-tests/all.lua ]; then \
      rm -rf lua-tests && \
      git clone -q -b go-lua https://github.com/Shopify/lua-tests.git lua-tests && \
      git -C lua-tests checkout -q 5ab7086f03ce05f12cb7c52cdb8df868cf8581be; \
    fi

RUN go mod download

RUN GOBIN=/usr/local/bin go install github.com/jstemmer/go-junit-report/v2@v2.1.0

RUN go build ./...

CMD ["/bin/bash"]
