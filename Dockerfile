from golang:latest
run mkdir /go/src/series_detector
copy hello.go /go/src/series_detector
run cd /go/src/series_detector && go build && go install
cmd series_detector 
