package sql

import "io"

type StringReader struct {
	buffer []rune
	pos    int
}

func (s *StringReader) Peek(length int) *string {
	if (s.pos + length) >= (len(s.buffer) - 1) {
		return nil
	}

	res := string(s.buffer[s.pos : s.pos+length])
	return &res
}

func (s *StringReader) PeekNext() rune {
	if s.pos >= len(s.buffer) {
		return -1
	}

	return s.buffer[s.pos]
}

func (s *StringReader) Next() rune {
	if s.pos >= len(s.buffer) {
		return -1
	}

	res := s.buffer[s.pos]
	s.pos++
	return res
}

func (s *StringReader) Consume(length int) error {
	_, err := s.Read(length)

	if err != nil {
		return err
	}
	return nil
}

func (s *StringReader) Read(length int) (*string, error) {
	res := s.Peek(length)
	if res != nil {
		s.pos += length
		return res, nil
	}

	return nil, io.EOF
}

func (s *StringReader) ReadString() string {
	next := s.PeekNext()
	buffer := make([]rune, 0)
	for next >= 'A' && next <= 'Z' || next >= 'a' && next <= 'z' {
		buffer = append(buffer, s.Next())
		next = s.PeekNext()
	}

	return string(buffer)
}

func (s *StringReader) PeekRemaining() string {
	return string(s.buffer[s.pos:])
}

func NewStringReader(sql *string) *StringReader {
	return &StringReader{
		buffer: []rune(*sql),
		pos:    0,
	}
}
