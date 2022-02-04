FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM
WORKDIR /build
RUN echo "$BUILDPLATFORM -> $TARGETPLATFORM" > /log.txt
COPY * .
RUN go build -o jmod-example .

FROM alpine
COPY --from=build /build/jmod-example .
COPY --from=build /build/webcomponent.js .
RUN mkdir -p /data

# EXPOSE 8080
CMD ["/jmod-example"]
