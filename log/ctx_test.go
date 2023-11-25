package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fioepq9/devtools/log"
)

var _ = Describe("Ctx", func() {
	var (
		buf     *bytes.Buffer
		l       *slog.Logger
		nowUnix = time.Now().Unix()
	)
	BeforeEach(func() {
		buf = &bytes.Buffer{}
		l = log.DefaultConsoleLogger(buf, slog.LevelDebug).With(slog.Int64("unix", nowUnix))
		slog.SetDefault(l)
	})

	Context("without logger", func() {
		It("should use default logger", func() {
			ll := log.FromContext(context.TODO())
			ll.Info("use default logger")
			Expect(buf.String()).To(ContainSubstring(fmt.Sprint(nowUnix)))
			Expect(buf.String()).To(ContainSubstring("use default logger"))
		})
	})

	Context("with logger", func() {
		It("should use context logger", func() {
			ctxNowUnix := time.Now().Unix()
			ctx := log.WithContext(context.TODO(), l.With(slog.Int64("ctx_unix", ctxNowUnix)))
			ll := log.FromContext(ctx)
			ll.Info("use context logger")
			Expect(buf.String()).To(ContainSubstring(fmt.Sprint(nowUnix)))
			Expect(buf.String()).To(ContainSubstring(fmt.Sprint(ctxNowUnix)))
			Expect(buf.String()).To(ContainSubstring("use context logger"))
		})
	})
})

var _ = Describe("Default", func() {
	var (
		cbuf *bytes.Buffer
		jbuf *bytes.Buffer
		cl   *slog.Logger
		jl   *slog.Logger
	)
	BeforeEach(func() {
		cbuf = &bytes.Buffer{}
		jbuf = &bytes.Buffer{}
		cl = log.DefaultConsoleLogger(cbuf, slog.LevelDebug)
		jl = log.DefaultJSONLogger(jbuf, slog.LevelDebug)
	})

	Context("use console logger", func() {
		It("should log", func() {
			cl.Info("console output ok")
			Expect(cbuf.String()).To(ContainSubstring("console output ok"))
		})
	})

	Context("use json logger", func() {
		It("should log", func() {
			jl.Info("json output ok")
			v := make(map[string]any)
			err := json.Unmarshal(jbuf.Bytes(), &v)
			Expect(err).To(BeNil())
			Expect(jbuf.String()).To(ContainSubstring("json output ok"))
		})
	})
})
