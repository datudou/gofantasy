package gofantasy

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

type marshaller interface {
	marshal(value any) ([]byte, error)
}

type jsonMarshaller struct{}

func (*jsonMarshaller) marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

type xmlMarshaller struct{}

func (*xmlMarshaller) marshal(value any) ([]byte, error) {
	return xml.Marshal(value)
}

type decoder interface {
	decode(reader io.Reader, into any) error
}

type xmlDecoder struct{}

func (*xmlDecoder) decode(reader io.Reader, into any) error {
	return xml.NewDecoder(reader).Decode(into)
}

type jsonDecoder struct{}

func (*jsonDecoder) jsonDecoder(reader io.Reader, into any) error {
	return json.NewDecoder(reader).Decode(into)
}
