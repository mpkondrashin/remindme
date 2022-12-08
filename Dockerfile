FROM scratch

COPY remindme /
COPY web /web

ENTRYPOINT [ "./remindme" ]
