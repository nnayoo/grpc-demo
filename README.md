

安装 go
PATH=$PATH:/usr/local/go/bin/

安装 

$ PB_REL="https://github.com/protocolbuffers/protobuf/releases"
$ curl -LO $PB_REL/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip

unzip protoc-3.15.8-linux-x86_64.zip -d $HOME/.local
export PATH="$PATH:$HOME/.local/bin"


$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct

$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

$ export PATH="$PATH:$(go env GOPATH)/bin"


protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/proto.proto



git config --global  user.name 'your name'

git config --global  user.email  'your email'

git init

git add

git commit 

git remote add origin git@github.com:tianqixin/runoob-git-test.git

git push -u origin main