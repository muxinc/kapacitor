FROM buildpack-deps:jessie-curl

ARG VERSION
RUN gpg \
    --keyserver hkp://ha.pool.sks-keyservers.net \
    --recv-keys 05CE15085FC09D18E99EFB22684A14CF2582E0C5

COPY kapacitor_${VERSION}_amd64.deb /
RUN dpkg -i /kapacitor_${VERSION}_amd64.deb
RUN rm -f /kapacitor_${VERSION}_amd64.deb

EXPOSE 9092

VOLUME /var/lib/kapacitor

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
CMD ["kapacitord"]
