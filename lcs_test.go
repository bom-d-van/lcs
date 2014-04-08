package lcs

import (
	"strings"

	"testing"
)

var fixtures = [][3]*stringProse{
	{newStringProse("abcdefg"), newStringProse("abckefg"), newStringProse("abcefg")},           // abc[d](k)efg
	{newStringProse("abc"), newStringProse("everything is abc or not"), newStringProse("abc")}, // (everything is )abc( or not)
	{newStringProse("nothing"), newStringProse("every e s abc r"), newStringProse("")},         // [nothing](every e s abc r)
}

var diffs = []string{
	"abc[d](k)efg",
	"(everything is )abc( or not)",
	"[nothing](every e s abc r)",
}

var longfixtures = [][]*stringProse{
	{newStringProse("abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg"), newStringProse("xxxxxxabckefg"), newStringProse("abcefg")}, // abc[d](k)efg
}

func TestLongLCS(t *testing.T) {
	for _, f := range longfixtures {
		exp := LCS(f[0], f[1])
		if exp.(*stringProse).String() != f[2].String() {
			t.Errorf("LCS(\"%s\", \"%s\"), expect: \"%s\", got: \"%s\"", f[0], f[1], f[2], exp)
		}
		// println(mylcs(f[0].String(), f[1].String()))
	}
}

func TestLCS(t *testing.T) {
	for _, f := range fixtures {
		exp := LCS(f[0], f[1])
		if exp.(*stringProse).String() != f[2].String() {
			t.Errorf("LCS(\"%s\", \"%s\"), expect: \"%s\", got: \"%s\"", f[0], f[1], f[2], exp)
		}
		// println(mylcs(f[0].String(), f[1].String()))
	}
}

func TestDiff(t *testing.T) {
	for i, f := range fixtures {
		exp := LCS(f[0], f[1])
		diff := Diff(f[0], f[1], exp)
		if diff.(*stringProse).String() != diffs[i] {
			t.Errorf("Diff(\"%s\", \"%s\"), expect: \"%s\", got: \"%s\"", f[0], f[1], diffs[i], diff)
		}
	}
}

var wordfixtures = [][]string{
	{
		"a simple word that is in test shape.a simple word that is in test shape.a simple word that is in test shape.a simple word that is in test shape.",
		"a simple that is in test shape.a simple that is in test shape.a simple that is in test shape.a simple that is in test shape.a simple that is in test shape.",
		"a simple [word] that is in test shape.",
	},
	// 	{
	// 		`he was an old man who fished alone in a skiff in the Gulf Stream and he had gone eighty-four days now without taking a fish. In the first forty days a boy had been with him. But after forty days without a fish the boy’s parents had told him that the old man was now definitely and finally salao, which is the worst form of unlucky, and the boy had gone at their orders in another boat which caught three good fish the first week. It made the boy sad to see the old man come in each day with his skiff empty and he always went down to help him carry either the coiled lines or the gaff and harpoon and the sail that was furled around the mast. The sila was patched with flour sacks and, furled, it looked like the flag of permanent defeat.

	// The old man was thin and gaunt with deep wrinkles in the back of his neck. The brown blotches of the benevolent skin cancer the sun brings from its [9] reflection on the tropic sea were on his cheeks. The blotches ran well down the sides of his face and his hands had the deep-creased scars from handling heavy fish on the cords. But none of these scars were fresh. They were as old as erosions in a fishless desert.

	// Everything about him was old except his eyes and they were the same color as the sea and were cheerful and undefeated.`,
	// 		`He was an old man who fished alone in a skiff in the Gulf Stream and he had gone eighty-four days now without taking a fish. In the first forty days a boy had been with him. But after forty days without a fish the boy’s parents had told him that the old man was now definitely and finally salao, which is the worst form of unlucky, and the boy had gone at their orders in another boat which caught three good fish the first week. It made the boy sad to see the old man come in each day with his skiff empty and he always went down to help him carry either the coiled lines or the gaff and harpoon and the sail that was furled around the mast. The sail was patched with flour sacks and, furled, it looked like the flag of permanent defeat.

	// The old man was thin and gaunt with deep wrinkles in the back of his neck. The brown blotches of the benevolent skin cancer the sun brings from its [9] reflection on the tropic sea were on his cheeks. The blotches ran well down the sides of his face and his hands had the deep-creased scars from handling heavy fish on the cords. But none of these scars were fresh. They were as old as erosions in a fishless desert.

	// Everything about him was old except his eyes and they were the same color as the sea and were cheerful and undefeated.`,

	// 		`[he](He) was an old man who fished alone in a skiff in the Gulf Stream and he had gone eighty-four days now without taking a fish. In the first forty days a boy had been with him. But after forty days without a fish the boy’s parents had told him that the old man was now definitely and finally salao, which is the worst form of unlucky, and the boy had gone at their orders in another boat which caught three good fish the first week. It made the boy sad to see the old man come in each day with his skiff empty and he always went down to help him carry either the coiled lines or the gaff and harpoon and the sail that was furled around the mast. The [sila](sail) was patched with flour sacks and, furled, it looked like the flag of permanent defeat.

	// The old man was thin and gaunt with deep wrinkles in the back of his neck. The brown blotches of the benevolent skin cancer the sun brings from its [9] reflection on the tropic sea were on his cheeks. The blotches ran well down the sides of his face and his hands had the deep-creased scars from handling heavy fish on the cords. But none of these scars were fresh. They were as old as erosions in a fishless desert.

	// Everything about him was old except his eyes and they were the same color as the sea and were cheerful and undefeated.`,
	// 	},
}

func TestArticle(t *testing.T) {
	for _, f := range wordfixtures {
		ori, edit := newArticle(strings.NewReader(f[0])), newArticle(strings.NewReader(f[1]))
		diff := Diff(ori, edit, LCS(ori, edit))
		if diff.(*article).String() != f[2] {
			t.Errorf("Diff(\"%s\", \"%s\"), expect: \"%s\", got: \"%s\"", f[0], f[1], f[2], diff)
		}
	}
}
