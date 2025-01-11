FROM --platform=linux/amd64 golang:alpine as builder 

# CREATE A FOLDER FOR OUR BUILD
RUN mkdir /build
ADD . /build/
WORKDIR /build


# GENERATE THE EXE BUILD FILE
RUN go mod tidy
RUN go build -o main .

# STAGE 2 
FROM alpine

# CREATE A USER
RUN adduser -S -D -H -h /app appuser

USER appuser

COPY . /app

COPY --from=builder /build/main /app/

WORKDIR /app

EXPOSE 8080

# RUN MAIN
CMD ["/app/main"]