package tui

import (
	"github.com/mi8bi/ghqx/internal/status"
)

// ViewState は TUI の表示状態を表す
type ViewState int

const (
	ViewStateList ViewState = iota // リスト表示
	ViewStateLoading                // ロード中
	ViewStateError                  // エラー表示
	ViewStateConfirm                // 確認ダイアログ
)

// Message はユーザー向けメッセージ
type Message struct {
	Text string
	Type MessageType
	Hint string // オプショナルなヒント
}

// MessageType はメッセージの種類
type MessageType int

const (
	MessageTypeInfo MessageType = iota
	MessageTypeSuccess
	MessageTypeWarning
	MessageTypeError
)

// OperationType は TUI で実行可能な操作
type OperationType int

const (
	OperationRefresh OperationType = iota
	OperationQuit
)

// ProjectRow はテーブル表示用のプロジェクト行
type ProjectRow struct {
	status.ProjectDisplay // 埋め込み
}

// NewProjectRow は ProjectDisplay から ProjectRow を作成する
func NewProjectRow(p status.ProjectDisplay) ProjectRow {
	row := ProjectRow{
		ProjectDisplay: p,
	}

	return row
}
