FROM golang:1.22.0

RUN useradd -m -s /bin/bash app && \
    apt update && \
    mv /usr/local/go /home/app/go

ENV GOPATH=/home/app/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOBIN

RUN chown -R -v app:app /home/app

USER app

WORKDIR /home/app/workspace

COPY --chown=app:app . .

CMD ["go", "run", "playground/main.go"]