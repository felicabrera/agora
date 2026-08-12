# ÁGORA API server image.
#
# Multi-stage, ending at distroless/static as a non-root user: the runtime layer
# has no shell, no package manager and no libc, so a remote-code-execution bug in
# the API has almost nothing to pivot to.
#
# The build stage keeps .git out (see .dockerignore) but the toolchain still
# stamps module information; -trimpath keeps absolute build paths out of the
# binary so the same source produces the same bytes on any machine.

FROM golang:1.26-bookworm AS build
WORKDIR /src

# Copy the module graph first so dependency download is cached independently of
# source edits.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/agora-api ./cmd/agora-api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/agora-api /usr/local/bin/agora-api
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/agora-api"]
