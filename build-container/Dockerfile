FROM alpine:3.12
COPY --chmod=555 prepare-import /app/prepare-import
WORKDIR /app
EXPOSE 3000
CMD ["/app/prepare-import"]