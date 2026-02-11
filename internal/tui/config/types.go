package configtui

import (
	"strconv"
	"sort" // Added this import

	"github.com/mi8bi/ghqx/internal/config"
)

// EditState は編集状態を表す
type EditState int

const (
	EditStateList EditState = iota // リスト表示
	EditStateEdit                   // 編集中
	EditStateSaving                 // 保存中
)

// Field は編集可能なフィールド
type Field struct {
	Name        string      // 表示名
	Key         string      // 内部キー
	Value       string      // 現在の値
	DefaultValue string     // デフォルト値
	Description string      // 説明
	Type        FieldType   // フィールドタイプ
	Options     []string    // FieldTypeSelection の場合に選択肢を保持
}

// FieldType はフィールドの種類
type FieldType int

const (
	FieldTypeString FieldType = iota
	FieldTypeBool
	FieldTypeInt
	FieldTypeSelection // Added this
)

// ConfigEditor は設定エディタのデータ
type ConfigEditor struct {
	Config      *config.Config
	Fields      []Field
	Modified    bool
	ConfigPath  string
}

// NewConfigEditor は ConfigEditor を作成する
func NewConfigEditor(cfg *config.Config, configPath string) *ConfigEditor {
	editor := &ConfigEditor{
		Config:     cfg,
		ConfigPath: configPath,
		Modified:   false,
	}

	editor.buildFields()
	return editor
}

// buildFields は編集可能なフィールドを構築する
func (e *ConfigEditor) buildFields() {
	// Get sorted root names for selection options
	var rootNames []string
	for name := range e.Config.Roots {
		rootNames = append(rootNames, name)
	}
	sort.Strings(rootNames)

	e.Fields = []Field{
		{
			Name:        "dev ルート",
			Key:         "roots.dev",
			Value:       e.Config.Roots["dev"],
			DefaultValue: "",
			Description: "dev ワークスペースのパス",
			Type:        FieldTypeString,
		},
		{
			Name:        "release ルート",
			Key:         "roots.release",
			Value:       e.Config.Roots["release"],
			DefaultValue: "",
			Description: "release ワークスペースのパス",
			Type:        FieldTypeString,
		},
		{
			Name:        "sandbox ルート",
			Key:         "roots.sandbox",
			Value:       e.Config.Roots["sandbox"],
			DefaultValue: "",
			Description: "sandbox ワークスペースのパス",
			Type:        FieldTypeString,
		},
		{
			Name:        "デフォルトルート",
			Key:         "default.root",
			Value:       e.Config.Default.Root,
			DefaultValue: "dev",
			Description: "デフォルトで使用するルート",
			Type:        FieldTypeSelection, // Changed to FieldTypeSelection
			Options:     rootNames,          // Populated with root names
		},
	}
}

// UpdateField はフィールドの値を更新する
func (e *ConfigEditor) UpdateField(index int, value string) {
	if index < 0 || index >= len(e.Fields) {
		return
	}

	e.Fields[index].Value = value
	e.Modified = true
}

// ApplyChanges は変更を Config に反映する
func (e *ConfigEditor) ApplyChanges() {
	for _, field := range e.Fields {
		switch field.Key {
		case "roots.dev":
			e.Config.Roots["dev"] = field.Value
		case "roots.release":
			e.Config.Roots["release"] = field.Value
		case "roots.sandbox":
			e.Config.Roots["sandbox"] = field.Value
		case "default.root":
			e.Config.Default.Root = field.Value
		}
	}
}

// ヘルパー関数

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func stringToBool(s string) bool {
	return s == "true" || s == "yes" || s == "1"
}

func intToString(i int) string {
	return strconv.Itoa(i)
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}