// Package encode provides base64 and hex encoding/decoding for log lines.
package encode

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// Encoding specifies the encoding scheme to apply.
type Encoding string

const (
	Base64  Encoding = "base64"
	Base64URL Encoding = "base64url"
	Hex     Encoding = "hex"
)

// Options configures the Encoder.
type Options struct {
	Encoding Encoding
	Decode   bool // if true, decode instead of encode
}

// Encoder encodes or decodes each input line.
type Encoder struct {
	opts Options
}

// New creates a new Encoder with the given options.
func New(opts Options) (*Encoder, error) {
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	return &Encoder{opts: opts}, nil
}

func validateOptions(opts Options) error {
	switch opts.Encoding {
	case Base64, Base64URL, Hex:
		return nil
	case "":
		return fmt.Errorf("encode: encoding scheme must be specified")
	default:
		return fmt.Errorf("encode: unknown encoding %q; want base64, base64url, or hex", opts.Encoding)
	}
}

// Apply encodes or decodes each line and returns the transformed lines.
func (e *Encoder) Apply(lines []string) ([]string, error) {
	out := make([]string, 0, len(lines))
	for i, line := range lines {
		transformed, err := e.transform(line)
		if err != nil {
			return nil, fmt.Errorf("encode: line %d: %w", i+1, err)
		}
		out = append(out, transformed)
	}
	return out, nil
}

func (e *Encoder) transform(line string) (string, error) {
	if e.opts.Decode {
		return e.decode(line)
	}
	return e.encode(line), nil
}

func (e *Encoder) encode(line string) string {
	switch e.opts.Encoding {
	case Base64:
		return base64.StdEncoding.EncodeToString([]byte(line))
	case Base64URL:
		return base64.URLEncoding.EncodeToString([]byte(line))
	case Hex:
		return hex.EncodeToString([]byte(line))
	default:
		return line
	}
}

func (e *Encoder) decode(line string) (string, error) {
	var b []byte
	var err error
	switch e.opts.Encoding {
	case Base64:
		b, err = base64.StdEncoding.DecodeString(line)
	case Base64URL:
		b, err = base64.URLEncoding.DecodeString(line)
	case Hex:
		b, err = hex.DecodeString(line)
	}
	if err != nil {
		return "", err
	}
	return string(b), nil
}
