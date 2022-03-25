FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY main /app/main  ##main.exe for windows

EXPOSE 9000

USER nonroot:nonroot

ENTRYPOINT ["/app/main"]