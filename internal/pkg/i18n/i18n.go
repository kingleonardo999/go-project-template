package i18n

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	mu      sync.RWMutex
	locales = map[string]map[string]string{}
)

// Load 加载指定语言的 locale 文件，lang 如 "zh-CN"、"en-US"
func Load(lang, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	m := make(map[string]string)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	mu.Lock()
	locales[lang] = m
	mu.Unlock()
	return nil
}

// T 返回指定语言中 code 对应的消息，找不到时返回 code 本身
func T(lang, code string) string {
	mu.RLock()
	m, ok := locales[lang]
	mu.RUnlock()
	if !ok {
		return code
	}
	if msg, ok := m[code]; ok {
		return msg
	}
	return code
}
