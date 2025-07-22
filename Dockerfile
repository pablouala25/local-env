FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap .

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /app/bootstrap /var/runtime/bootstrap
CMD [ "bootstrap" ]
