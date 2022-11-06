FROM alpine
COPY client_golang /app/client_golang
WORKDIR /app 
CMD ["./client_golang"]