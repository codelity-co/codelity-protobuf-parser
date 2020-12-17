package parser

import (
	"errors"
)

type Rpc struct {
	name string
	request string
	response string
}

func (r *Rpc) GetRpcName() string {
	return r.name
}

func (r *Rpc) GetRpcRequestName() string {
	return r.request
}

func (r *Rpc) GetRpcResponseName() string {
	return r.response
}

type Service struct {
	name string
	rpcs []*Rpc
}

func (s *Service) addRpc(rpc *Rpc) error {
	if s.rpcs == nil {
		return errors.New("s.rpcs is nil")
	}
	s.rpcs = append(s.rpcs, rpc)
  return nil
}

func (s *Service) GetServiceName() string {
  return s.name
}

func (s *Service) GetRpcs() []*Rpc {
  return s.rpcs
}