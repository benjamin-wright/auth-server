FROM busybox

COPY app /app
COPY www/ /www/

ENTRYPOINT [ "/app" ]