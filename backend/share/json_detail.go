package share

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// JsonToStruct json を 渡された struct で解析して struct に詰める
// goStruct には ポインタで渡すこと
func JsonToStruct(goStruct interface{}, path string) error {
	// ファイルの有無を調べる
	_, err := os.Stat(path)
	// ファイルが存在するか
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("ファイルが存在しません: %w", err)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// JSON から struct にする
	if err := json.Unmarshal(bytes, &goStruct); err != nil {
		return err
	}
	// fmt.Printf("%#v\n", goStruct)
	return nil
}

// saveJSON JSON ファイルに保存
// bytes は JSON として保存したい 構造体のバイト
// 本来は infra 層 に 全部書くべきと思われる
func SaveJSON(bytes []byte, path string) error {
	// ファイルを開く(write 権限付き)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 全部入り配列を一気に書き込む
	if _, err := file.Write(bytes); err != nil {
		return err
	}

	return nil
}

// パスと保存したい内容をまとめた構造体
type syncJSON struct {
	path  string
	bytes []byte
}

func newSyncJSON(bytes []byte, path string) (syncJSON, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return syncJSON{}, fmt.Errorf("ファイルが存在しません: %w", err)
	}
	if errors.Is(err, fs.ErrNotExist) {
		return syncJSON{}, fmt.Errorf("ファイルが存在しません: %w", err)
	}

	return syncJSON{bytes: bytes, path: path}, nil
}

// old ファイル名を作成
func (s *syncJSON) oldFile() string {
	filename := s.path
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	return name + "_old" + ext
}

// 同期して一気に保存したいデータたち
type syncJSONs []syncJSON

func NewSyncJSONs(bytes []byte, path string) (syncJSONs, error) {
	newJSON, err := newSyncJSON(bytes, path)
	if err != nil {
		return nil, err
	}
	s := syncJSONs{}
	return append(s, newJSON), nil
}

func (ss *syncJSONs) Add(bytes []byte, path string) error {
	addJSON, err := newSyncJSON(bytes, path)
	if err != nil {
		return err
	}
	*ss = append(*ss, addJSON)
	return nil
}

func fileCopy(baseFile string, newFile string) error {
	// ファイルを開く
	src, err := os.Open(baseFile)
	if err != nil {
		return err
	}
	defer src.Close()

	// 新しいファイルを作成する
	dst, err := os.Create(newFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	// ファイルをコピーする
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

// 同期しながら保存したい対象ファイルの old を作る
func (ss *syncJSONs) toOld() error {
	for _, v := range *ss {
		err := fileCopy(v.path, v.oldFile())
		if err != nil {
			return err
		}
	}
	return nil
}

// old ファイルを削除する
func (ss *syncJSONs) rmOld() error {
	for _, v := range *ss {
		err := os.Remove(v.oldFile())
		if err != nil {
			return err
		}
	}
	return nil
}

// old ファイルから戻す
func (ss *syncJSONs) bkOld() error {
	for _, v := range *ss {
		// 書き込み途中のファイルを削除する
		err := os.Remove(v.path)
		if err != nil {
			return err
		}
		// old ファイルから戻す
		err02 := fileCopy(v.oldFile(), v.path)
		if err02 != nil {
			return err02
		}
	}
	return nil
}

// SaveSyncJSON ロールバックを考慮して すべて保存 or すべて保存しない で動く
// 保存する前のファイルを一時的にコピーしておいて 現ファイルを編集するけど
// エラーが発生したら 一時ファイルで現ファイルを上書きする
func (ss *syncJSONs) Save() error {
	err_toOld := ss.toOld()
	if err_toOld != nil {
		return err_toOld
	}

	var err error
	for _, v := range *ss {
		err = SaveJSON(v.bytes, v.path)
		if err != nil {
			break
		}
	}

	if err != nil {
		// バックアップから戻す
		err_bkOld := ss.bkOld()
		if err_bkOld != nil {
			return fmt.Errorf("%w, 元のエラーは %v", err_bkOld, err)
		}
	}

	err_rmOld := ss.rmOld()
	if err_rmOld != nil {
		if err != nil {
			return fmt.Errorf("%w, 元のエラーは %v", err_rmOld, err)
		}
		return err_rmOld
	}
	return nil
}
