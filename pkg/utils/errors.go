package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
)

func MentorError(msg string, code codes.Code, details protoiface.MessageV1) error {
	s := status.New(code, msg)
	ns, err := s.WithDetails(details)
	if err != nil {
		return err
	}
	return ns.Err()
}
