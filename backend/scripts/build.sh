PROJ=github.com/foundcenter/moas/backend
SRC_DIR=$GOPATH/src/$PROJ
OUT_DIR=$SRC_DIR/bin/
BIN_NAME=moas
rm $OUT_DIR/$BIN_NAME
mkdir -p $OUT_DIR
go build -o $OUT_DIR/$BIN_NAME $PROJ/
go fmt $PROJ
go vet $PROJ