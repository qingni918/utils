package netutil

import (
	"errors"
	"net"
	"os"
	"strings"
	"syscall"
)

// ErrorType 表示网络错误类型
type ErrorType int

const (
	UnknownError ErrorType = iota
	ConnectionRefused
	Timeout
	ConnectionReset
	Temporary
	EOFError
)

func DetectNetErrorType(err error) ErrorType {
	if err == nil {
		return UnknownError
	}

	// 1. Connection Refused
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			if sysErr.Err == syscall.ECONNREFUSED {
				return ConnectionRefused
			}
			if sysErr.Err == syscall.ECONNRESET {
				return ConnectionReset
			}
		}

		// 2. 超时错误
		if opErr.Timeout() {
			return Timeout
		}

		// 3. 临时性错误
		if opErr.Temporary() {
			return Temporary
		}
	}

	// 4. 判断超时接口
	if nerr, ok := err.(net.Error); ok {
		if nerr.Timeout() {
			return Timeout
		}
		//if nerr.Temporary() {
		//	return Temporary
		//}
	}

	// 5. EOF 错误
	if errors.Is(err, os.ErrClosed) || strings.Contains(err.Error(), "EOF") {
		return EOFError
	}

	return UnknownError
}

func IsConnectionRefused(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			return sysErr.Err == syscall.ECONNREFUSED
		}
	}
	return false
}
