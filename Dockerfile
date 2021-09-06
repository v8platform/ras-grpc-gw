FROM scratch
EXPOSE 3002
COPY ras-grpc-gw /
ENTRYPOINT ["/ras-grpc-gw"]