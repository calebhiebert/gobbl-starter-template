# Start with a golang image
FROM golang:1.10.3-stretch as build

# Install dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep && go get -u github.com/gobuffalo/packr/v2/packr2

# Create a user to run the app as
RUN useradd --shell /bin/bash bot

# Set the workdir to the application path
WORKDIR $GOPATH/src/bot

# Copy all application files
COPY . .

# Run the packr
RUN packr2

# Install packages
RUN dep ensure --vendor-only

# Build the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -ldflags="-w -s" -o /go/bin/bot

# Start from a scratch container for a nice and small image
FROM scratch

# Copy the binary build
COPY --from=build /go/bin/bot /go/bin/bot

# Copy the password file (with the bot user) from the build container
COPY --from=build /etc/passwd /etc/passwd

# Copy assets
COPY ./assets /go/bin/assets

# Set the assets directory
ENV ASSETS_DIRECTORY "/go/bin/assets"

# Set the user to the previously created user
USER bot

# Expose the API port
EXPOSE 8080

CMD [ "/go/bin/bot" ]