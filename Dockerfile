FROM node:18-alpine AS client-build
WORKDIR /app
COPY client/. .
RUN npm ci && npm run build

FROM golang:1.19 as api-build
WORKDIR /app
COPY backend/. ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/api


FROM alpine:latest 
WORKDIR /app
COPY --from=client-build /app/dist/client/* ./dist/
COPY --from=api-build /go/bin/api ./
CMD ./api