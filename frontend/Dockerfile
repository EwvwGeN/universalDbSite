FROM alpine
ADD main /server/
ADD ./configs/ /configs/
ADD ./web /web/
ARG ISCNF
ENV ISCNF=${ISCNF}
CMD ["sh", "-c", "./server/main $ISCNF"]