package dtls

import (
	"encoding/binary"
)

// https://www.iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml
type extensionValue uint16

const (
	extensionSupportedGroupsValue extensionValue = 10
	extensionUseSRTPValue                        = 14
)

type extension interface {
	marshal() ([]byte, error)
	unmarshal(data []byte) error

	extensionValue() extensionValue
}

func decodeExtensions(buf []byte) ([]extension, error) {
	declaredLen := binary.BigEndian.Uint16(buf)
	if len(buf)-2 != int(declaredLen) {
		return nil, errLengthMismatch
	}

	extensions := []extension{}
	unmarshalAndAppend := func(data []byte, e extension) error {
		err := e.unmarshal(data)
		if err != nil {
			return err
		}
		extensions = append(extensions, e)
		return nil
	}

	for offset := 2; offset+2 < len(buf); {
		var err error
		switch extensionValue(binary.BigEndian.Uint16(buf[offset:])) {
		case extensionSupportedGroupsValue:
			err = unmarshalAndAppend(buf[offset:], &extensionSupportedGroups{})
		default:
		}
		if err != nil {
			return nil, err
		}

		extensionLength := binary.BigEndian.Uint16(buf[offset+2:])
		offset += (2 + int(extensionLength))
	}

	return extensions, nil
}

func encodeExtensions(e []extension) ([]byte, error) {
	extensions := []byte{}
	for _, e := range e {
		raw, err := e.marshal()
		if err != nil {
			return nil, err
		}
		extensions = append(extensions, raw...)
	}
	out := []byte{0x00, 0x00}
	binary.BigEndian.PutUint16(out, uint16(len(extensions)))
	return append(out, extensions...), nil
}
