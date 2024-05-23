package main

import (
	"fmt"
	"os"
	"regexp"
	"unicode/utf8"
)

type Tile struct {
	Number rune
	Suit   rune
}

func (t Tile) String() string {
	return fmt.Sprintf("%c%c", t.Number, t.Suit)
}

func parseHaifu(haifu string) []Tile {
	var tiles []Tile
	re := regexp.MustCompile(`([1-9]+[psm])|([東南西北中発白])`)
	matches := re.FindAllStringSubmatch(haifu, -1)

	for _, match := range matches {
		if match[2] != "" { // 字牌の場合
			tileRune, _ := utf8.DecodeRuneInString(match[2])
			tiles = append(tiles, Tile{Number: tileRune, Suit: 'z'})
		} else if match[1] != "" { // 数牌の場合
			nums := match[1][:len(match[1])-1]
			suit := rune(match[1][len(match[1])-1])
			for _, num := range nums {
				tiles = append(tiles, Tile{Number: rune(num), Suit: suit})
			}
		}
	}

	// デバッグ用出力
	for _, tile := range tiles {
		fmt.Printf("解析された牌: %c%c\n", tile.Number, tile.Suit)
	}

	return tiles
}

func parseTumo(tumo string) Tile {
	re := regexp.MustCompile(`([1-9東南西北中発白][psmz]?)`)
	match := re.FindStringSubmatch(tumo)
	if len(match) > 0 {
		tileRune, _ := utf8.DecodeRuneInString(match[1])
		suit := 'z'
		if len(match[1]) > 1 {
			suit = rune(match[1][1])
		}
		return Tile{Number: tileRune, Suit: suit}
	}
	return Tile{}
}

func isWinningHand(tiles []Tile) bool {
	if len(tiles) != 14 {
		fmt.Println("牌の数が14枚ではありません。")
		return false
	}
	if isChiitoitsu(tiles) {
		return true
	}
	return isStandardWinningHand(tiles)
}

func isChiitoitsu(tiles []Tile) bool {
	tileCount := make(map[Tile]int)
	for _, tile := range tiles {
		tileCount[tile]++
	}

	pairCount := 0
	for _, count := range tileCount {
		if count == 2 {
			pairCount++
		} else if count != 2 {
			return false
		}
	}

	return pairCount == 7
}

func isStandardWinningHand(tiles []Tile) bool {
	tileCount := make(map[Tile]int)
	for _, tile := range tiles {
		tileCount[tile]++
	}

	fmt.Println("牌のカウント:", tileCount)

	for tile, count := range tileCount {
		if count >= 2 {
			newTileCount := copyTileCount(tileCount)
			newTileCount[tile] -= 2
			if newTileCount[tile] == 0 {
				delete(newTileCount, tile)
			}
			fmt.Println("対子を見つけました:", tile)
			fmt.Println("対子を除いた後の牌のカウント:", newTileCount)
			if checkMentsu(newTileCount) {
				return true
			}
		}
	}
	return false
}

func checkMentsu(tileCount map[Tile]int) bool {
	if len(tileCount) == 0 {
		return true
	}

	fmt.Println("面子チェック開始:", tileCount)

	for tile, count := range tileCount {
		if count >= 3 {
			newTileCount := copyTileCount(tileCount)
			newTileCount[tile] -= 3
			if newTileCount[tile] == 0 {
				delete(newTileCount, tile)
			}
			fmt.Println("刻子を見つけました:", tile)
			fmt.Println("刻子を除いた後の牌のカウント:", newTileCount)
			if checkMentsu(newTileCount) {
				return true
			}
		}

		if tile.Suit != 'z' {
			next1 := Tile{Number: tile.Number + 1, Suit: tile.Suit}
			next2 := Tile{Number: tile.Number + 2, Suit: tile.Suit}
			if tileCount[next1] > 0 && tileCount[next2] > 0 {
				newTileCount := copyTileCount(tileCount)
				newTileCount[tile]--
				newTileCount[next1]--
				newTileCount[next2]--
				if newTileCount[tile] == 0 {
					delete(newTileCount, tile)
				}
				if newTileCount[next1] == 0 {
					delete(newTileCount, next1)
				}
				if newTileCount[next2] == 0 {
					delete(newTileCount, next2)
				}
				fmt.Println("順子を見つけました:", tile, next1, next2)
				fmt.Println("順子を除いた後の牌のカウント:", newTileCount)
				if checkMentsu(newTileCount) {
					return true
				}
			}
		}
	}

	return false
}

func copyTileCount(tileCount map[Tile]int) map[Tile]int {
	newTileCount := make(map[Tile]int)
	for k, v := range tileCount {
		newTileCount[k] = v
	}
	return newTileCount
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <牌譜文字列> <ツモ牌>")
		return
	}

	haifu := os.Args[1]
	tumo := os.Args[2]

	parsedTiles := parseHaifu(haifu)
	tumoTile := parseTumo(tumo)

	parsedTiles = append(parsedTiles, tumoTile)

	// デバッグ用に牌の数を表示
	fmt.Println("入力された牌の数:", len(parsedTiles))

	if isWinningHand(parsedTiles) {
		fmt.Println("上がりです！")
	} else {
		fmt.Println("まだ上がりではありません。")
	}
}
