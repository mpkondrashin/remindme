FROM busybox:glibc
EXPOSE 80/tcp
WORKDIR /
COPY remindme /
COPY web /web
CMD [ "/remindme" ]
