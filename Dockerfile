FROM busybox:glibc
WORKDIR /
COPY remindme /
COPY web /web
CMD [ "/remindme" ]
