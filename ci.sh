git commit -am "x"
git push origin master
GOOS=linux GOARCH=amd64 go build -o record_exporter main.go && \
    docker build -t  ianleto/record_exporter:$(git rev-parse --short HEAD) -f Dockerfile .&&\
    docker push  ianleto/record_exporter:$(git rev-parse --short HEAD)