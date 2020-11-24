FROM ubuntu:bionic

# replace shell with bash so we can source files
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

# update the repository sources list
# and install dependencies
RUN apt-get update \
    && apt-get install -y curl gnupg2 \
    && apt-get -y autoclean

# nvm environment variables
ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION 12.19.0

# install nvm
# https://github.com/creationix/nvm#install-script
RUN curl --silent -o- https://raw.githubusercontent.com/creationix/nvm/v0.31.2/install.sh | bash

# install node and npm
RUN source $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

# add node and npm to path so the commands are available
ENV NODE_PATH $NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

# install yarn
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install yarn -y

# download and setup go
RUN curl https://dl.google.com/go/go1.15.3.linux-amd64.tar.gz -O
RUN tar -C /usr/local -xzf go1.15.3.linux-amd64.tar.gz
ENV PATH /usr/local/go/bin:$PATH

# confirm installations
RUN node -v
RUN npm -v
RUN yarn --version
RUN go version

# copy all files over
RUN mkdir /app
WORKDIR /app
COPY main.go .
COPY package.json .
COPY public ./public
COPY src ./src

# build client and server
RUN yarn install
RUN yarn build
RUN go build main.go

# show output (run with --progress plain to see)
RUN echo "showing build output" && ls

# set default command
CMD /app/main
