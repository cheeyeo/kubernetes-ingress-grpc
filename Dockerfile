FROM scratch
ADD certs/tls.crt /certs/tls.crt
ADD certs/tls.key /certs/tls.key
ADD server /
CMD ["/server"]
