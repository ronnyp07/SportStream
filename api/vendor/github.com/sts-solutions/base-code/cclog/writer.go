package cclog

import "bytes"

type CaptureWriter struct {
	Buf bytes.Buffer
}

func (cw *CaptureWriter) Write(p []byte) (n int, err error) {
	cw.Buf.Write(p)

	return len(p), nil
}

func (cw *CaptureWriter) Sync() error {
	return nil
}

func (cw *CaptureWriter) String() string {
	return cw.Buf.String()
}
