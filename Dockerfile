FROM alpine
COPY ./app/native /app/native
WORKDIR /app 
ENTRYPOINT ["./native"]