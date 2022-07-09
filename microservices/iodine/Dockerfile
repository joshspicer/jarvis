FROM ubuntu:20.04

RUN apt update -y \
    && apt install -y openssh-server git wget \
    && useradd -ms /bin/bash sshproxy \
    && echo "sshproxy:sshproxy" | chpasswd
    #  && echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config

RUN cd / \
    && git clone https://github.com/joshspicer/iodine.git

RUN cd /iodine \
    && make

RUN cd /iodine \
    && make install

RUN chmod +x /iodine/start.sh

ENTRYPOINT ["sh"]
CMD ["iodined"]