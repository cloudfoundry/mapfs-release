FROM harbor-repo.vmware.com/dockerhub-proxy-cache/relintdockerhubpushbot/cf-deployment-concourse-tasks as golang_version
RUN git clone --recurse-submodules https://github.com/cloudfoundry/mapfs-release
RUN cd mapfs-release && bosh create-release --tarball /tmp/release.tgz
RUN version=$(cat /tmp/release.tgz | tar -Oxz packages/golang-1-linux.tgz | tar z --list | grep -ohE "go[0-9]\.[0-9]{1,2}\.[0-9]{1,2}") && echo $version > /tmp/golang_version

FROM harbor-repo.vmware.com/dockerhub-proxy-cache/apnar/nfs-ganesha

RUN DEBIAN_FRONTEND=noninteractive \
     && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 3FE869A9 \
     && echo "deb http://ppa.launchpad.net/gluster/nfs-ganesha-2.5/ubuntu xenial main" > /etc/apt/sources.list.d/nfs-ganesha-2.5.list \
     && echo "deb http://ppa.launchpad.net/gluster/libntirpc-1.5/ubuntu xenial main" > /etc/apt/sources.list.d/libntirpc-1.5.list \
     && echo "deb http://ppa.launchpad.net/gluster/glusterfs-3.13/ubuntu xenial main" > /etc/apt/sources.list.d/glusterfs-3.13.list \
     && apt-get update \
     && apt-get install -y nfs-ganesha-mem curl

RUN apt-get update
RUN apt-get install -y git
RUN apt-get install -y iproute
RUN apt-get install -y fuse
RUN apt-get install -y build-essential
RUN apt-get install -y wget

COPY --from=golang_version /tmp/golang_version /tmp/golang_version
RUN wget https://dl.google.com/go/$(cat /tmp/golang_version).linux-amd64.tar.gz
RUN tar -xvf $(cat /tmp/golang_version).linux-amd64.tar.gz
RUN mv go /usr/local

ENV PATH "$PATH:/usr/local/go/bin"

RUN sed -i 's/VFS/MEM/' /start.sh
RUN sed -i 's#^EXPORT#NFSV4 \{ Graceless = true; \}\nEXPORT#' /start.sh
