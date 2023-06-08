package encoding

import (
	"bytes"
	"io"
	"strings"
)

func Encode(encoder Encoder, writer io.Writer, reader io.Reader) error {
	return encoder.Encode(writer, reader)
}

func EncodeString(encoder Encoder, str string) ([]byte, error) {
	writer := &bytes.Buffer{}
	err := Encode(encoder, writer, strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func Decode(decoder Decoder, writer io.Writer, reader io.Reader) error {
	return decoder.Decode(writer, reader)
}

func DecodeString(decoder Decoder, reader io.Reader) (string, error) {
	writer := &strings.Builder{}
	err := Decode(decoder, writer, reader)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}
