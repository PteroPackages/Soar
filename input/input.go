package input

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func Parse(input string) (map[string]string, error) {
	m := map[string]string{}
	r := strings.NewReader(input)

	for {
		if err := readSpace(r); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		k, err := readKey(r)
		if err != nil {
			return nil, err
		}
		fmt.Println(k)

		if err := readSpace(r); err != nil {
			if errors.Is(err, io.EOF) {
				return nil, errors.New("missing value for key")
			}

			return nil, err
		}

		v, err := readValue(r)
		if err != nil {
			return nil, err
		}
		fmt.Println(v)

		m[k] = v
	}

	return m, nil
}

func readSpace(r *strings.Reader) error {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}

		// 		\/ byte(' ')
		if b != 32 {
			return nil
		}
	}
}

func readKey(r *strings.Reader) (string, error) {
	if err := r.UnreadByte(); err != nil {
		return "", err
	}

	b := strings.Builder{}

	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}

		if c == '=' {
			break
		}

		b.WriteRune(c)
	}

	return b.String(), nil
}

func readValue(r *strings.Reader) (string, error) {
	if err := r.UnreadByte(); err != nil {
		return "", err
	}

	b := strings.Builder{}
	c, _, err := r.ReadRune()
	if err != nil {
		return "", err
	}

	d := c == '"'
	if !d {
		b.WriteRune(c)
	}

	for {
		c, _, err = r.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return "", err
		}

		if (d && c == '"') || c == ' ' {
			break
		}

		b.WriteRune(c)
	}

	return b.String(), nil
}
