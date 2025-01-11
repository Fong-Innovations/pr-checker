FROM --platform=linux/amd64 golang:latest as builder 


RUN go get -u github.com/cosmtrek/air

# CREATE A FOLDER FOR OUR BUILD
RUN mkdir /build
ADD . /build/
WORKDIR /build


# GENERATE THE EXE BUILD FILE
RUN go mod tidy
RUN go build -o main .

# STAGE 2 
FROM alpine

COPY . /app

COPY --from=builder /build/main /app/

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["air"]