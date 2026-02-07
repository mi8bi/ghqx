package tui

import "github.com/mi8bi/ghqx/internal/domain"

// ViewState は TUI の表示状態を表す
type ViewState int

const (
	ViewStateList    ViewState = iota // リスト表示
	ViewStateLoading                  // ロード中
	ViewStateError                    // エラー表示
	ViewStateConfirm                  // 確認ダイアログ
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
	OperationPromote OperationType = iota
	OperationUndo
	OperationRefresh
	OperationQuit
)

// ProjectRow はテーブル表示用のプロジェクト行
type ProjectRow struct {
	Project     domain.Project
	CanPromote  bool   // プロモート可能か
	PromoteHint string // プロモートできない理由
}

// NewProjectRow は ProjectRow を作成する
func NewProjectRow(project domain.Project) ProjectRow {
	row := ProjectRow{
		Project:    project,
		CanPromote: false,
	}

	// sandbox のプロジェクトのみプロモート可能
	if project.Zone == domain.ZoneSandbox {
		if project.Dirty {
			row.CanPromote = false
			row.PromoteHint = "変更をコミットしてください"
		} else {
			row.CanPromote = true
		}
	} else {
		row.PromoteHint = "sandbox のプロジェクトのみプロモート可能"
	}

	return row
}
