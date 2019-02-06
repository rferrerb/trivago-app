FROM scratch
ADD main /
ADD config.toml /
ADD Index.tmpl /
CMD ["/main"]
