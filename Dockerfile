FROM scratch
WORKDIR /
COPY remindme /
COPY web /web
CMD [ /remindme ]
