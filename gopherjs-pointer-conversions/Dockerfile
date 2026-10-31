FROM public.ecr.aws/d3j8x8q7/olympus-base-go:latest

WORKDIR /app
COPY . .

# GopherJS at this commit only supports Go 1.21: it refuses to run with any
# other Go version, and it compiles the standard library of its GOROOT to
# JavaScript, so the GOROOT must be a Go 1.21 distribution. The base image
# ships a newer Go, so the matching release is installed next to it
# (checksum-verified) and put first on PATH.
RUN curl -fsSL https://go.dev/dl/go1.21.13.linux-amd64.tar.gz -o /tmp/go.tar.gz \
    && echo "502fc16d5910562461e6a6631fb6377de2322aad7304bf2bcd23500ba9dab4a7  /tmp/go.tar.gz" | sha256sum -c - \
    && mkdir -p /usr/local/go1.21 \
    && tar -C /usr/local/go1.21 --strip-components=1 -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz

ENV PATH=/usr/local/go1.21/bin:/root/go/bin:$PATH
ENV GOTOOLCHAIN=local
ENV GOFLAGS=-buildvcs=false

RUN go mod download

RUN GOBIN=/usr/local/bin go install github.com/jstemmer/go-junit-report/v2@v2.1.0

# GopherJS runs the compiled tests with Node.js. The base image already ships
# Node.js (Debian's nodejs package, /usr/bin/node), so nothing is installed
# here; the build fails early if it is ever missing.
RUN test -x /usr/bin/node && /usr/bin/node --version

# `gopherjs test` and `gopherjs run` pass node a --stack_size derived from the
# process stack limit. When that limit is unlimited, as in some sandboxes, the
# derived size is enormous and node exits with status 255 before running
# anything. This wrapper (ahead of /usr/bin/node on PATH) uses node's usual
# 8 MiB stack in that case and leaves every other invocation unchanged.
RUN printf '%s\n' \
      '#!/bin/bash' \
      'if [ "$(ulimit -s)" = unlimited ]; then' \
      '  ulimit -s 8192 2>/dev/null' \
      '  set -- "${@/#--stack_size=*/--stack_size=8128}"' \
      'fi' \
      'exec /usr/bin/node "$@"' \
      > /usr/local/bin/node \
    && chmod 0755 /usr/local/bin/node

RUN go install .

CMD ["/bin/bash"]
