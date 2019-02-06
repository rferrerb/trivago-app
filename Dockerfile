FROM scratch
ADD main /
ADD config.toml /
CMD ["/main"]
