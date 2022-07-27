FROM gcr.io/distroless/static-debian11:debug
ENTRYPOINT ["/protog"]
COPY protog /