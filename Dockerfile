FROM golang:1.16-buster AS BUILD

RUN apt update && apt install -y libwebp-dev && mkdir /project

WORKDIR /project

COPY ./ ./

RUN go build

FROM ubuntu AS RUNTIME

RUN ln -fs /usr/share/zoneinfo/Europe/Amsterdam /etc/localtime \
    && apt update \
    && apt install -y poppler-utils libwebp6 libreoffice \
    fonts-dejavu fonts-freefont-ttf fonts-ubuntu ttf-bitstream-vera

COPY --from=BUILD /project/document-preview-microservice /

EXPOSE 3030

CMD [ "sh", "-c", "/document-preview-microservice" ]
