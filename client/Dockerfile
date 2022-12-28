FROM node:13.12.0-alpine

# set work dir
WORKDIR /app

# set env
ENV PATH /app/node_modules/.bin:$PATH
ENV PORT=80

# set dependencies
COPY package.json /app
RUN npm install --silent

# copy over src
COPY . /app

# start the app
ENTRYPOINT [ "/app/entrypoint.sh" ]
