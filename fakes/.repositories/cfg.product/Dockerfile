FROM ubuntu:21.04 AS builder
LABEL stage=autodelete

RUN apt update && apt install -y curl gnupg software-properties-common apt-transport-https
RUN curl https://baltocdn.com/helm/signing.asc | apt-key add -
RUN echo "deb https://baltocdn.com/helm/stable/debian/ all main" | tee /etc/apt/sources.list.d/helm-stable-debian.list
RUN apt update && apt install helm

RUN curl -LO "https://dl.k8s.io/release/v1.23.3/bin/linux/amd64/kubectl"
RUN curl -LO "https://dl.k8s.io/v1.23.3/bin/linux/amd64/kubectl.sha256"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

RUN curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -
RUN apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
RUN apt update && apt install terraform

# -----------------------------------------------

FROM golang:1.17 as mage
LABEL stage=autodelete

RUN git clone https://github.com/magefile/mage
RUN cd mage && go run bootstrap.go
RUN mv $GOPATH/bin/mage /bin/mage

# -----------------------------------------------

FROM docker.io/aryazanov/prodctl:0.0.1 AS prodctl
LABEL stage=autodelete

# -----------------------------------------------

FROM docker.io/aryazanov/bundle:0.0.1 AS bundle
LABEL stage=autodelete

# -----------------------------------------------

FROM golang:1.17.7 AS runtime

COPY --from=builder /bin/helm /usr/local/bin/helm
COPY --from=builder /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=builder /bin/terraform /usr/local/bin/terraform
COPY --from=mage /bin/mage /usr/local/bin/mage
COPY --from=prodctl /bin/prodctl /bundle/prodctl
COPY --from=bundle /bundle/ /bundle/

ENV PATH="/bundle:${PATH}"

WORKDIR /bundle

USER root