FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM
WORKDIR /build
RUN echo "$BUILDPLATFORM -> $TARGETPLATFORM" > /log.txt
COPY * .
RUN go build -o testmod .

FROM alpine
COPY --from=build /build/* .
RUN mkdir -p /data

# EXPOSE 8080
CMD ["/testmod"]
