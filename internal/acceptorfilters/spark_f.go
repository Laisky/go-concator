package acceptorfilters

import (
	"fmt"
	"regexp"

	"gofluentd/library"
	"gofluentd/library/log"

	"github.com/Laisky/zap"
)

type SparkFilterCfg struct {
	IgnoreRegex                   *regexp.Regexp
	Name, MsgKey, Tag, Identifier string
}

// SparkFilter filter spark messages.
// some old spark messages need tobe discard
type SparkFilter struct {
	*BaseFilter
	*SparkFilterCfg
}

func NewSparkFilter(cfg *SparkFilterCfg) *SparkFilter {
	log.Logger.Info("NewSparkFilter",
		zap.String("regex", cfg.IgnoreRegex.String()),
		zap.String("tag", cfg.Tag))

	if cfg.Identifier == "" {
		panic(fmt.Errorf("`Identifier` should not be empty, but got: %v", cfg.Identifier))
	}

	return &SparkFilter{
		BaseFilter:     &BaseFilter{},
		SparkFilterCfg: cfg,
	}
}

func (f *SparkFilter) GetName() string {
	return f.Name
}

func (f *SparkFilter) Filter(msg *library.FluentMsg) *library.FluentMsg {
	if msg.Tag != f.Tag {
		return msg
	}

	switch msg.Message[f.MsgKey].(type) {
	case []byte:
	default:
		return msg
	}

	// discard some format
	// log.Logger.Debug("ignore spark log",
	// 	zap.String("tag", f.Tag),
	// 	zap.ByteString("log", msg.Message[f.MsgKey].([]byte)))
	if f.IgnoreRegex.Match(msg.Message[f.MsgKey].([]byte)) {
		f.DiscardMsg(msg)
		return nil
	}

	// set spark container_id
	msg.Message[f.Identifier] = []byte("spark")

	return msg
}
