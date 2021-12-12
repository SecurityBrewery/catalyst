FROM ubuntu:18.04

RUN apt-get update -y && apt-get -y install curl gnupg2 software-properties-common
RUN curl -OL https://download.arangodb.com/arangodb34/DEBIAN/Release.key
RUN apt-key add Release.key
RUN apt-add-repository 'deb https://download.arangodb.com/arangodb34/DEBIAN/ /'
RUN apt-get update -y && apt-get -y install arangodb3

COPY catalyst /app/catalyst
CMD /app/catalyst

EXPOSE 8000
