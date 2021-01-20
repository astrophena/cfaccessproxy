FROM gcr.io/distroless/static
COPY cfaccessproxy /
CMD ["/cfaccessproxy"]
