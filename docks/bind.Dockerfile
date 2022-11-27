FROM docker.uclv.cu/ubuntu:bionic

RUN apt-get update \
  && apt-get install -y \
  iputils-ping \
  bind9 \
  bind9utils \
  bind9-doc

# Enable IPv4
RUN sed -i 's/OPTIONS=.*/OPTIONS="-4 -u bind"/' /etc/default/bind9

# Cnfiguration files
VOLUME [ "/etc/bind/" ]
# VOLUME [ "/etc/bind/named.conf.options" ]
# Domain files
VOLUME [ "/var/cache/bind" ]

# Run eternal loop
CMD ["/etc/init.d/bind9", "start", "&&", "/bin/bash", "-c", "while :; do sleep 10; done"]
