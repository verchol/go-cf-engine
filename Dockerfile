FROM alpine
COPY ./dist /dist
CMD ["/dist/linux_386"]