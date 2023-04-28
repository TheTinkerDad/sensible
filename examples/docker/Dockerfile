FROM nginxdemos/hello

COPY sensible sensible
COPY wrapper.sh wrapper.sh
RUN chmod +x /wrapper.sh

EXPOSE 8090

CMD . ./wrapper.sh