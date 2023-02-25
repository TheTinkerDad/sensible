FROM nginxdemos/hello

COPY sensible sensible
COPY wrapper.sh wrapper.sh
COPY scripts/ /etc/sensible/

RUN chmod +x /wrapper.sh
RUN chmod +x /etc/sensible/*.sh

EXPOSE 8090

CMD . ./wrapper.sh