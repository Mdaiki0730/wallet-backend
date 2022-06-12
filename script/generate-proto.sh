set -e

PROTOBUF_PATH=./api/proto
PROTO_THIRD_PARTY=${PROTOBUF_PATH}/third_party
SERVICES=`ls -F ${PROTOBUF_PATH}|grep /|grep -v "third_party"|sed 's/\/$//'`
for service in ${SERVICES}; do
  rm -rf ${PROTOBUF_PATH}/${service}/${service}pb
  mkdir -p ${PROTOBUF_PATH}/${service}/${service}pb
  exists=`ls ${PROTOBUF_PATH}/${service}/*.proto >/dev/null 2>&1; echo $?`
  if [ ! $exists -ne 0 ]; then
    protoc \
    -I ${PROTO_THIRD_PARTY} \
    --proto_path=${PROTOBUF_PATH}/${service} ${PROTOBUF_PATH}/${service}/*.proto \
    --go-grpc_out=require_unimplemented_servers=false:${PROTOBUF_PATH}/${service}/${service}pb \
    --go_out=${PROTOBUF_PATH}/${service}/${service}pb \
    --grpc-gateway_out=${PROTOBUF_PATH}/${service}/${service}pb;
  fi
done
