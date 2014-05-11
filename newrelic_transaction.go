package newrelic_go_api

/*
#cgo CFLAGS: -fopenmp -I./_include
#cgo LDFLAGS: -lnewrelic-collector-client
#cgo LDFLAGS: -lnewrelic-common
#cgo LDFLAGS: -lnewrelic-transaction

#include <stdlib.h> 
#include <newrelic_transaction.h> 
#include <newrelic_collector_client.h> 
static void setupEmbededCollectorCGOProxy() {
    newrelic_register_message_handler(newrelic_message_handler);
}

*/
import "C"

import (
	"unsafe"
)

type TTransactionId int64
type SegmentId int64

const (
	NR_ERROR_CODE_OK                      = 0
	NR_ERROR_CODE_INVALID_ID              = 1
	NR_ERROR_CODE_TRANSACTION_NOT_STARTED = 2
	NR_ERROR_CODE_TRANSACTION_IN_PROGRESS = 3
	NR_ERROR_CODE_TRANSACTION_NOT_NAMED   = 4

	NR_DATASTORE_OPERATION_SELECT = "select"
	NR_DATASTORE_OPERATION_INSERT = "insert"
	NR_DATASTORE_OPERATION_UPDATE = "update"
	NR_DATASTORE_OPERATION_DELETE = "delete"
)

func StartWebTransaction() TTransactionId {
	result := C.newrelic_transaction_begin()
	return TTransactionId(result)
}

func NameWebTransaction(transactionId TTransactionId, name string) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	result := C.newrelic_transaction_set_name(C.long(transactionId), cName)
	return int(result)
}

func EndWebTransaction(transactionId TTransactionId) int {
	result := C.newrelic_transaction_end(C.long(transactionId))
	return int(result)
}

func StartDatastoreStatement(tId TTransactionId, sId SegmentId, table string, operation string) SegmentId {
	if operation != NR_DATASTORE_OPERATION_SELECT &&
		operation != NR_DATASTORE_OPERATION_INSERT &&
		operation != NR_DATASTORE_OPERATION_UPDATE &&
		operation != NR_DATASTORE_OPERATION_DELETE {
		return SegmentId(0)
	}

	cTable := C.CString(table)
	defer C.free(unsafe.Pointer(cTable))

	cOperation := C.CString(operation)
	defer C.free(unsafe.Pointer(cOperation))

	cId := C.long(tId)
	cSeg := C.long(sId)

	result := C.newrelic_segment_datastore_begin(cId, cSeg, cTable, cOperation)
	return SegmentId(result)
}

func EndDatastoreStatement(transactionId TTransactionId, segmentId SegmentId) int {
	result := C.newrelic_segment_end(C.long(transactionId), C.long(segmentId))
	return int(result)
}

/// Only call once and after nr_init
func SetupEmbeddedCollectorClient() {
	C.setupEmbededCollectorCGOProxy()
}
