FROM alpine:latest
ADD redis-cl /
ADD redis-conf.json /
EXPOSE 20221
ENTRYPOINT ["/redis-cl", "--logtostderr"]