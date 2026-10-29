FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -o /out/taskapi ./cmd/taskapi

FROM scratch
COPY --from=build /out/taskapi /taskapi
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/taskapi", "-addr", ":8080"]
