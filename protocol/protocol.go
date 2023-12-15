package protocol

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"laser/constants"
)

type FileMessage struct {
	FileNameSize int
	FileName     string
	ContentSize  uint64
	Content      []byte
}

func (m *FileMessage) WriteTo(writer io.Writer) (int64, error) {
	var written int64
	startByte := byte(constants.StartPattern)
	n, err := writer.Write([]byte{startByte})
	if err != nil {
		return written, err
	}
	written += int64(n)

	// FileNameSize를 varint로 변환하여 작성
	buf := make([]byte, binary.MaxVarintLen64)
	n = binary.PutUvarint(buf, uint64(m.FileNameSize))
	if _, err := writer.Write(buf[:n]); err != nil {
		return written, err
	}
	written += int64(n)

	// FileName 작성
	n, err = writer.Write([]byte(m.FileName))
	if err != nil {
		return written, err
	}
	written += int64(n)

	// ContentSize를 varint로 변환하여 작성
	n = binary.PutUvarint(buf, m.ContentSize)
	if _, err := writer.Write(buf[:n]); err != nil {
		return written, err
	}
	written += int64(n)

	// Content 작성
	n, err = writer.Write(m.Content)
	if err != nil {
		return written, err
	}
	written += int64(n)

	return written, nil
}

func ReadFileMessage(r io.Reader) (*FileMessage, error) {
	defer fmt.Println("read done")
	reader := bufio.NewReader(r)
	_, err := reader.ReadBytes(constants.StartPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to read start pattern: %v", err)
	}
	var msg FileMessage
	fmt.Println("read start")

	// FileNameSize를 varint로 읽음
	fileNameSize, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil, err
	}
	msg.FileNameSize = int(fileNameSize)

	// FileName 읽음
	fileName := make([]byte, fileNameSize)
	if _, err := io.ReadFull(reader, fileName); err != nil {
		return nil, err
	}
	msg.FileName = string(fileName)

	// ContentSize를 varint로 읽음
	contentSize, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil, err
	}
	msg.ContentSize = contentSize

	// Content 읽음
	content := make([]byte, contentSize)
	if _, err := io.ReadFull(reader, content); err != nil {
		return nil, err
	}
	msg.Content = content

	return &msg, nil
}
