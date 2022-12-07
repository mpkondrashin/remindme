FROM scratch

COPY remindme /

ENTRYPOINT [ "./remindme" ]
