# Pull the Node image from Docker Hub
FROM node:17-alpine AS build

# Setting Working Directory
WORKDIR /usr/app

# set env
ENV PATH /app/node_modules/.bin:$PATH

# Copying only package.json
COPY package*.json ./

# Install Dependencies
RUN npm install --silent
RUN npm install react-scripts@5.0.1 -g --silent

# Copy rest of the code to container
COPY . .

EXPOSE 3000

# start the app
# ENTRYPOINT [ "/app/entrypoint.sh" ]
CMD ["npm", "start"]