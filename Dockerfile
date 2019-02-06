FROM scratch
ADD main /
ADD config.toml /
ADD *.tmpl /
CMD ["/main"]
