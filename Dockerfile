FROM golang:1.6-onbuild
ENTRYPOINT ["/api"]
ADD api /
