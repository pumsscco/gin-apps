package main

import (
	"testing"
)

// benchmarking the decode function
func BenchmarkCB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCombo()
	}
}

func BenchmarkMI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getMission()
	}
}

func BenchmarkRM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getRoleMan()
	}
}

func BenchmarkET(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getEnemyType("男性")
		getEnemyType("女性")
		getEnemyType("其它")
	}
}

func BenchmarkFT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFighterType("内功心法", "neigong")
		getFighterType("情侣合技", "lover")
		getFighterType("怒技", "rage")
		getFighterType("普通招式", "common")
		getFighterType("组合技", "combo")
	}
}

func BenchmarkIT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getItemType("扇")
		getItemType("剑")
		getItemType("短剑")
		getItemType("弓")

		getItemType("盔甲")
		getItemType("鞋")
		getItemType("佩饰")

		getItemType("武功")

		getItemType("丹药")
		getItemType("暗器")
		getItemType("食物")
		
		getItemType("食材")
	}
}