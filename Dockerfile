FROM scratch

COPY remindme /
COPY web /web

CMD [ "/remindme" ]
