package libs

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"bytes"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestMarshalUnmarshalFluentBatchMsg(t *testing.T) {
	v := FluentBatchMsg{}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgFluentBatchMsg(b *testing.B) {
	var (
		v   = FluentBatchMsg{}
		err error
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = v.MarshalMsg(nil); err != nil {
			b.Fatalf("%+v", err)
		}
	}
}

func BenchmarkAppendMsgFluentBatchMsg(b *testing.B) {
	var (
		v   = FluentBatchMsg{}
		err error
	)
	bts := make([]byte, 0, v.Msgsize())
	if bts, err = v.MarshalMsg(bts[0:0]); err != nil {
		b.Fatalf("%+v", err)
	}
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if bts, err = v.MarshalMsg(bts[0:0]); err != nil {
			b.Fatalf("%+v", err)
		}
	}
}

func BenchmarkUnmarshalFluentBatchMsg(b *testing.B) {
	v := FluentBatchMsg{}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeFluentBatchMsg(t *testing.T) {
	var (
		v   = FluentBatchMsg{}
		err error
	)
	var buf bytes.Buffer
	if err = msgp.Encode(&buf, &v); err != nil {
		t.Fatalf("%+v", err)
	}

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := FluentBatchMsg{}
	err = msgp.Decode(&buf, &vn)
	if err != nil {
		t.Error(err)
	}

	buf.Reset()
	if err = msgp.Encode(&buf, &v); err != nil {
		t.Fatalf("%+v", err)
	}
	err = msgp.NewReader(&buf).Skip()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeFluentBatchMsg(b *testing.B) {
	var (
		err error
		v   = FluentBatchMsg{}
		buf bytes.Buffer
	)
	if err = msgp.Encode(&buf, &v); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = v.EncodeMsg(en); err != nil {
			b.Fatalf("%+v", err)
		}
	}
	en.Flush()
}

func BenchmarkDecodeFluentBatchMsg(b *testing.B) {
	var (
		v   = FluentBatchMsg{}
		err error
		buf bytes.Buffer
	)
	if err = msgp.Encode(&buf, &v); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = v.DecodeMsg(dc); err != nil {
			b.Fatal(err)
		}
	}
}

func TestMarshalUnmarshalFluentMsg(t *testing.T) {
	v := FluentMsg{}
	bts, err := v.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	left, err := v.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after UnmarshalMsg(): %q", len(left), left)
	}

	left, err = msgp.Skip(bts)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left over after Skip(): %q", len(left), left)
	}
}

func BenchmarkMarshalMsgFluentMsg(b *testing.B) {
	var (
		v   = FluentMsg{}
		err error
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = v.MarshalMsg(nil); err != nil {
			b.Fatalf("marshal: %+v", err)
		}
	}
}

func BenchmarkAppendMsgFluentMsg(b *testing.B) {
	v := FluentMsg{}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalMsg(bts[0:0])
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bts, _ = v.MarshalMsg(bts[0:0])
	}
}

func BenchmarkUnmarshalFluentMsg(b *testing.B) {
	v := FluentMsg{}
	bts, _ := v.MarshalMsg(nil)
	b.ReportAllocs()
	b.SetBytes(int64(len(bts)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestEncodeDecodeFluentMsg(t *testing.T) {
	var (
		buf bytes.Buffer
		v   = FluentMsg{}
		err error
	)
	if err = msgp.Encode(&buf, &v); err != nil {
		t.Fatal(err)
	}

	m := v.Msgsize()
	if buf.Len() > m {
		t.Logf("WARNING: Msgsize() for %v is inaccurate", v)
	}

	vn := FluentMsg{}
	if err = msgp.Decode(&buf, &vn); err != nil {
		t.Error(err)
	}

	buf.Reset()
	if err = msgp.Encode(&buf, &v); err != nil {
		t.Fatal(err)
	}
	if err = msgp.NewReader(&buf).Skip(); err != nil {
		t.Error(err)
	}
}

func BenchmarkEncodeFluentMsg(b *testing.B) {
	var (
		v   = FluentMsg{}
		err error
		buf bytes.Buffer
	)
	if err = msgp.Encode(&buf, &v); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(buf.Len()))
	en := msgp.NewWriter(msgp.Nowhere)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = v.EncodeMsg(en); err != nil {
			b.Fatalf("%+v", err)
		}
	}
	en.Flush()
}

func BenchmarkDecodeFluentMsg(b *testing.B) {
	var (
		v   = FluentMsg{}
		err error
		buf bytes.Buffer
	)
	if err = msgp.Encode(&buf, &v); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(buf.Len()))
	rd := msgp.NewEndlessReader(buf.Bytes(), b)
	dc := msgp.NewReader(rd)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = v.DecodeMsg(dc); err != nil {
			b.Fatal(err)
		}
	}
}
