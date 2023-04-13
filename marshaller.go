package gofantasy

import "encoding/json"

type marshaller interface {
	marshal(value any) ([]byte, error)
}

type jsonMarshaller struct{}

func (jm *jsonMarshaller) marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

type xmlMarshaller struct{}

func (xml *xmlMarshaller) marshal(value any) ([]byte, error) {
	return xml.marshal(value)
}
