FROM golang:alpine

WORKDIR ${GOROOT}/bda

ADD go.mod .
ADD go.sum .
RUN go mod download && go mod verify

ADD api                 ./api
ADD connection          ./connection
ADD crawler             ./crawler
ADD crawler_collector   ./crawler_collector
ADD logger              ./logger
ADD models              ./models
ADD pinger              ./pinger
ADD pinger_collector    ./pinger_collector
ADD types               ./types
ADD utils               ./utils
ADD *.go                ./

RUN go build -v -o ./nodes
RUN mv ./nodes /bin

WORKDIR /

# setup cron 2 hour job tasks
ADD ./crawler_collector/crawler_cron.sh /
RUN chmod +x /crawler_cron.sh
#RUN echo "0     */2       *       *       *       run-parts /etc/periodic/2hours" >> /etc/crontabs/root
RUN echo "*/1     *       *       *       *       /crawler_cron.sh > /dev/stdout" >> /etc/crontabs/root
RUN crontab -l

# set number of file descriptors over 16k - 16k max stable connections
RUN ulimit -n 16384
