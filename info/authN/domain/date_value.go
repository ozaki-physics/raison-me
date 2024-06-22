package domain

import (
	"time"
)

// 日付 値オブジェクト
type Date struct {
	time.Time
}

func constructorDate(d time.Time) (Date, DomainError) {
	if d.IsZero() {
		return Date{time.Now()}, NewDomainError("日付を入力してください")
	}

	return Date{d}, nil
}

func NewDate() (Date, DomainError) {
	d := time.Now()
	return constructorDate(d)
}

func ReNewDate(d string) (Date, DomainError) {
	dt, err := time.Parse(time.RFC3339, d)
	if err != nil {
		return Date{time.Now()}, WrapDomainError("日付のパースに失敗しました", err)
	}
	return constructorDate(dt)
}

// 以下ゲッター

// このサービス内で使う日付表記のフォーマット
func (d *Date) MyFormat() string {
	return d.Format(time.RFC3339)
}
