package library

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// LoadTagsAppendEnv append env to tags slice
//
// []string{"spring"} -> []string{"spring.sit"}
//
// Deprecated: should not auto append env, use LoadTagsReplaceEnv replaced
func LoadTagsAppendEnv(env string, tags []string) []string {
	ret := []string{}
	for _, t := range tags {
		ret = append(ret, t+"."+env)
	}

	return ret
}

// LoadTagReplaceEnv replace `{env}` to env
//
// "spring.{env}" -> "spring.sit"
func LoadTagReplaceEnv(env, tag string) string {
	return strings.ReplaceAll(tag, "{env}", env)
}

// LoadTagsReplaceEnv replace `{env}` to env
//
// []string{"spring.{env}"} -> []string{"spring.sit"}
func LoadTagsReplaceEnv(env string, tags []string) []string {
	nts := make([]string, len(tags))
	for i, t := range tags {
		nts[i] = strings.ReplaceAll(t, "{env}", env)
	}
	return nts
}

// LoadTagsMapAppendEnv append env to tags map
//
// map[string]interface{}{"spring": "xxx"} -> map[string]interface{}{"spring.sit": "xxx"}
//
// // Deprecated: should not auto append env
func LoadTagsMapAppendEnv(env string, tags map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	for t, v := range tags {
		ret[t+"."+env] = v
	}

	return ret
}

func RegexNamedSubMatch(r *regexp.Regexp, log []byte, subMatchMap map[string]interface{}) error {
	matches := r.FindSubmatch(log)
	names := r.SubexpNames()
	if len(names) != len(matches) {
		return errors.New("the number of args in `regexp` and `str` not matched")
	}

	for i, name := range names {
		if i != 0 && name != "" && len(matches[i]) != 0 {
			subMatchMap[name] = bytes.TrimSpace(matches[i])
		}
	}
	return nil
}

func FlattenMap(data map[string]interface{}, delimiter string) {
	for k, vi := range data {
		if v2i, ok := vi.(map[string]interface{}); ok {
			FlattenMap(v2i, delimiter)
			for k3, v3i := range v2i {
				data[k+delimiter+k3] = v3i
			}
			delete(data, k)
		}
	}
}

var defaultTemplateWithMappReg = regexp.MustCompile(`(?sm)\$\{([^}]+)\}`)

// TemplateWithMap replace `${var}` in template string
func TemplateWithMap(tpl string, data map[string]interface{}) string {
	return TemplateWithMapAndRegexp(defaultTemplateWithMappReg, tpl, data)
}

// TemplateWithMapAndRegexp replace `${var}` in template string
func TemplateWithMapAndRegexp(tplReg *regexp.Regexp, tpl string, data map[string]interface{}) string {
	var (
		k, vs string
	)
	for _, kg := range tplReg.FindAllStringSubmatch(tpl, -1) {
		k = kg[1]
		switch v := data[k].(type) {
		case string:
			vs = v
		case []byte:
			vs = string(v)
		case int:
			vs = strconv.FormatInt(int64(v), 10)
		case int64:
			vs = strconv.FormatInt(v, 10)
		case float64:
			vs = strconv.FormatFloat(v, 'f', -1, 64)
		}
		tpl = strings.ReplaceAll(tpl, fmt.Sprintf("${%v}", k), vs)
	}

	return tpl
}

// AbsInt return abs(int)
func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
