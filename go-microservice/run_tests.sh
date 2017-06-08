
CURRENT_DIR=$(dirname ${BASH_SOURCE[0]})
pushd $CURRENT_DIR > /dev/null
CURRENT_DIR=$(pwd)

export GOPATH=$CURRENT_DIR

export CXX="/opt/irods-externals/clang3.8-0/bin/clang++"
export C="/opt/irods-externals/clang3.8-0/bin/clang"

go get github.com/jjacquay712/GoRODS

go test msibasic_example/msibasic_example_test.go msibasic_example/msibasic_example.go -v