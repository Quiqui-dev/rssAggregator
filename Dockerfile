FROM golang:1.22-bullseye as build_stage

WORKDIR /app

COPY go.mod go.sum ./

COPY . .
RUN go build -o /rssAgg


# COPY --from=build_stage /rssAgg /rssAgg

EXPOSE 8080

CMD [ "/rssAgg"]