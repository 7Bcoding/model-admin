package grpc

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
)

var _ error = &KratosError{}

const (
	ErrorNotFound              = "NOT_FOUND_ERROR"
	ErrorUsing                 = "USING_ERROR"
	ErrorNotSupport            = "NOT_SUPPORT_ERROR"
	ErrorInvalidParams         = "INVALID_PARAMS_ERROR"
	ErrorWatchTimeout          = "WATCH_TIMEOUT_ERROR"
	ErrorInsufficientResource  = "INSUFFICIENT_RESOURCE_ERROR"
	ErrorInvalidObjectType     = "INVALID_OBJECT_TYPE_ERROR"
	ErrorServerUnavailable     = "SERVER_UNAVAILABLE_ERROR"
	ErrorConflict              = "CONFLICT_ERROR"
	ErrorScheduleUnAvailable   = "SCHEDULE_UNAVAILABLE_ERROR"
	ErrorVolumeUnavailable     = "VOLUME_UNAVAILABLE_ERROR"
	ErrorInvalidOperate        = "INVALID_OPERATE_ERROR"
	ErrorInvalidTxn            = "INVALID_TXN_ERROR"
	ErrorInvalidResourceChange = "INVALID_RESOURCE_CHANGE_ERROR"
	ErrorInvalidNotification   = "INVALID_NOTIFICATION_ERROR"
	ErrorNetworkNotReady       = "NETWORK_NOT_READY"
)

type KratosError struct {
	err *errors.Error
}

func (kerr *KratosError) Error() string {
	return kerr.err.Message
}

func NewKratosError(err *errors.Error) *KratosError {
	return &KratosError{err: err}
}

func NotFound(msg string) error {
	return NewKratosError(errors.New(400, ErrorNotFound, msg))
}

func IsNotFound(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorNotFound
	}
	return false
}

func Using(msg string) error {
	return NewKratosError(errors.New(400, ErrorUsing, msg))
}

func IsUsing(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorUsing
	}
	return false
}

func NotSupport(msg string) error {
	return NewKratosError(errors.New(400, ErrorNotSupport, msg))
}

func IsNotSupport(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorNotSupport
	}
	return false
}

func InvalidParams(msg string) error {
	return NewKratosError(errors.New(400, ErrorInvalidParams, msg))
}

func IsInvalidParams(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorInvalidParams
	}
	return false
}

func WatchTimeout(resource string) error {
	return NewKratosError(errors.New(400, ErrorWatchTimeout, fmt.Sprintf("watch %s timeout.", resource)))
}

func InsufficientResource(msg string) error {
	return NewKratosError(errors.New(400, ErrorInsufficientResource, msg))
}

func IsInsufficientResource(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorInsufficientResource
	}
	return false
}

func InvalidObjectType(resource string) error {
	return NewKratosError(errors.New(400, ErrorInvalidObjectType, fmt.Sprintf("invalid object %s type.", resource)))
}

func IsInvalidObjectType(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorInvalidObjectType
	}
	return false
}

func InvalidFilterParams(msg string) error {
	return NewKratosError(errors.New(400, ErrorInvalidObjectType, msg))
}

func ServerUnavailable(server string) error {
	return NewKratosError(errors.New(400, ErrorServerUnavailable, fmt.Sprintf("server %s status error. please check server.", server)))
}

func Conflict(msg string) error {
	return NewKratosError(errors.New(400, ErrorConflict, msg))
}

func IsConflict(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorConflict
	}
	return false
}

func ScheduleUnAvailable(msg string) error {
	return NewKratosError(errors.New(400, ErrorScheduleUnAvailable, msg))
}

func IsScheduleUnAvailable(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorScheduleUnAvailable
	}
	return false
}

func VolumeUnAvailable(msg string) error {
	return NewKratosError(errors.New(400, ErrorVolumeUnavailable, msg))
}

func IsVolumeUnAvailable(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorVolumeUnavailable
	}
	return false
}

func InvalidOperate(msg string) error {
	return NewKratosError(errors.New(400, ErrorInvalidOperate, msg))
}

func IsInvalidOperate(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorInvalidOperate
	}
	return false
}

func TxnInvalid(msg string) error {
	return NewKratosError(errors.New(400, ErrorInvalidTxn, msg))
}

func IsTxnInvalid(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorInvalidTxn
	}
	return false
}

func InvalidResourceChange(msg string) error {
	return NewKratosError(errors.New(400, ErrorInvalidResourceChange, msg))
}

func IsInvalidResourceChange(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorInvalidResourceChange
	}
	return false
}

func InvalidNotification(msg string) error {
	return NewKratosError(errors.New(400, ErrorInvalidNotification, msg))
}

func IsInvalidNotification(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorInvalidNotification
	}
	return false
}

func NetworkNotReady(msg string) error {
	return NewKratosError(errors.New(400, ErrorNetworkNotReady, msg))
}

func IsNetworkNotReady(err error) bool {
	ke, ok := err.(*KratosError)
	if ok {
		return errors.Reason(ke.err) == ErrorNetworkNotReady
	}
	return false
}
