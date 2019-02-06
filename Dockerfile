FROM scratch
ADD main /
ENV database_user=$database_user
ENV database_password=$database_password
CMD ["/main"]
