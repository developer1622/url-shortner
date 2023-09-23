# todo: make it multistage

# base linux flavour with deps installed :)
# actually base image
FROM golang:1.21.1

# just set some workloads
WORKDIR /app

# copy package managers 
COPY go.mod go.sum ./

# let's download all deps
RUN go mod download

# copy everything 
COPY . .

# some deps, will see in depth why we need to do those things
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

# so that we can access it from our host system
EXPOSE 8087

# Run
CMD ["/main"]