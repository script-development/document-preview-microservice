FROM golang:1.19-bullseye AS BUILD

RUN apt-get update && apt-get install -y libwebp-dev && mkdir /project

WORKDIR /project

COPY ./ ./

RUN go build

FROM ubuntu:focal AS RUNTIME

ARG DEBIAN_FRONTEND=noninteractive

RUN ln -fs /usr/share/zoneinfo/Europe/Amsterdam /etc/localtime \
    && apt-get update \
    && apt-get install -y --no-install-recommends poppler-utils libwebp6 libreoffice \
    openjdk-17-jre-headless fonts-dejavu fonts-freefont-ttf fonts-ubuntu ttf-bitstream-vera \
    && update-ca-certificates \
    && apt-get autoremove \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY --from=BUILD /project/document-preview-microservice /

EXPOSE 3030

CMD [ "sh", "-c", "/document-preview-microservice" ]
