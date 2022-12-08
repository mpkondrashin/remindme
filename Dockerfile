FROM ubuntu
WORKDIR /
COPY remindme /
COPY web /web
CMD [ "/remindme" ]
