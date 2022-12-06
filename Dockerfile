FROM scratch

COPY ./remaindme .

ENTRYPOINT [ "./remindme" ]