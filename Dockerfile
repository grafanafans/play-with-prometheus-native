FROM alpine
COPY ./app/client_golang /app/client_golang
WORKDIR /app 
ENTRYPOINT ["./client_golang"]