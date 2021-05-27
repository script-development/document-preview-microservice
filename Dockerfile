FROM golang:1.16-alpine AS BUILD

RUN apk add --no-cache musl-dev gcc libwebp-dev && mkdir /project

WORKDIR /project

COPY ./ ./

RUN go build

FROM alpine AS RUNTIME

RUN apk add --no-cache poppler-utils libwebp \
    && apk add ttf-dejavu font-noto ttf-ubuntu-font-family \
    msttcorefonts-installer \
    && update-ms-fonts

COPY --from=BUILD /project/document-preview-microservice /

EXPOSE 3030

CMD [ "sh", "-c", "/document-preview-microservice" ]
