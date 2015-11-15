FROM alpine:3.2
ADD explorer-srv /explorer-srv
ENTRYPOINT [ "/explorer-srv" ]
