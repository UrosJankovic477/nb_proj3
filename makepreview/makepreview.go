package makepreview

import (
	"bytes"

	"github.com/tcolgate/mp3"
)

func ProcessSong(data []byte) ([]byte, uint16, error) {
	buf := bytes.NewBuffer(data)
	mp3dec := mp3.NewDecoder(buf)
	length := 0.0
	previewbytes := 0
	for {
		frame := mp3.Frame{}
		skipped := 0
		err := mp3dec.Decode(&frame, &skipped)
		length += float64(frame.Duration().Seconds())
		if err != nil {
			break
		}
		if length <= 30 {
			previewbytes += frame.Size()
		}

	}

	return append([]byte(nil), data[:previewbytes]...), uint16(length), nil
}
