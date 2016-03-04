FROM alpine:3.3

RUN apk add --update git && rm -rf /var/cache/apk/*

EXPOSE 8123

ENTRYPOINT ["data-models-service", "-port", "8123", "-host", "0.0.0.0"]

CMD ["-log", "error", "-path", "/opt/repos", "-repo", "https://github.com/chop-dbhi/data-models"]

COPY data-models-service /usr/local/bin/
