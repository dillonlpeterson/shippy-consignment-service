# We use the official golang image, which contains all the 
# xorrect build tools and libraries. Notice 'as builder',
# this gives the container a name that we can reference later on.
FROM golang:1.9.0 as builder 

# Set our workdir to our current service in the gopath
WORKDIR /go/src/github.com/dillonlpeterson/shippy-consignment-service

# Copy the current directory into our workdir 
COPY . .

# Here we are pulling in godep, which is a dependency management tool
# We're going to use godep instead of go get (Gets fussy with sub-packages sometimes)
RUN go get -u github.com/golang/dep/cmd/dep

# Create a dep project, and run 'ensure', which will pull in all of the 
# dependencies of this directory.
RUN dep init && dep ensure 

# Build the binary, with a few flags that will allow us 
# to run the binary in Alpine. 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

# Here we're using a second FROM statement, which tells Docker to start 
# a new build process with this image
FROM alpine:latest 

# Security-related package good to have 
RUN apk --no-cache add ca-certificates

RUN mkdir /app 
WORKDIR /app 

#ADD consignment-service /app/consignment-service 
# Instead of pulling binary from Host machine, we pill it from the container named builder!
# Copies into /app directory, which is the current working directory.
COPY --from=builder /go/src/github.com/dillonlpeterson/shippy-consignment-service/consignment-service .

# As usual, run the binary!
CMD ["./consignment-service"]

# Code must be pushed up to Git so that it can pull in other services.