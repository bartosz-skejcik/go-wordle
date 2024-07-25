package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

var answers [][]Answer

type Game struct {
	word  string
	try   int
	ended bool
	won   bool
}

type Answer struct {
	idx      int
	letter   string
	correct  bool
	isInWord bool
}

func getRandomWordsList(list *[]string) error {
	if res, err := http.Get("https://random-word-api.vercel.app/api?words=10&length=5"); err != nil {
		return err
	} else {
		if body, err := io.ReadAll(res.Body); err != nil {
			return err
		} else {
			if err := json.Unmarshal(body, list); err != nil {
				return err
			}
		}
	}

	return nil
}

func randomRange(max, min int) int {
	return rand.IntN(max-min) + min
}

func clearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "linux", "darwin": // "darwin" is for macOS
		cmd = exec.Command("clear")
	default:
		fmt.Println("Unsupported platform")
		return
	}

	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("Error clearing the console:", err)
	}
}

func (g *Game) Init() {
	var availableWords []string
	randomIdx := randomRange(4, 1)
	if err := getRandomWordsList(&availableWords); err != nil {
		fmt.Printf("\n%s", err.Error())
		os.Exit(1)
	}
	g.try = 1
	g.word = availableWords[randomIdx-1]
}

func (g *Game) EndRound(guess string) {
	if guess == g.word {
		g.ended = true
		g.won = true
		return
	} else if g.try >= 5 {
		g.ended = true
		return
	}

	var ans []Answer
	guessArr := strings.Split(guess, "")
	wordArr := strings.Split(g.word, "")

	for i, _ := range guessArr {
		if guessArr[i] == wordArr[i] {
			ans = append(ans, Answer{letter: guessArr[i], correct: true, idx: i})
		} else if slices.Contains(wordArr, guessArr[i]) {
			ans = append(ans, Answer{letter: guessArr[i], correct: false, isInWord: true, idx: i})
		} else {
			ans = append(ans, Answer{letter: guessArr[i], correct: false, isInWord: false, idx: i})
		}
	}

	g.try += 1
	answers = append(answers, ans)
	g.ShowAnswer()

}

func (g *Game) ShowAnswer() {
	var t string

	for _, answer := range answers {
		var s string
		for _, item := range answer {
			if item.correct {
				// is correctly placed, green
				s += color.GreenString("%s ", item.letter)
			} else if !item.correct && item.isInWord {
				// is in word, color blue
				s += color.BlueString("%s ", item.letter)
			} else {
				s += fmt.Sprintf("%s ", item.letter)
			}
		}
		t += fmt.Sprintf("%v\n", s)
	}

	fmt.Println(t)

}

func main() {
	var game Game

	game.Init()

	fmt.Println(`
 __       __                            __  __
/  |  _  /  |                          /  |/  |
$$ | / \ $$ |  ______    ______    ____$$ |$$ |  ______
$$ |/$  \$$ | /      \  /      \  /    $$ |$$ | /      \
$$ /$$$  $$ |/$$$$$$  |/$$$$$$  |/$$$$$$$ |$$ |/$$$$$$  |
$$ $$/$$ $$ |$$ |  $$ |$$ |  $$/ $$ |  $$ |$$ |$$    $$ |
$$$$/  $$$$ |$$ \__$$ |$$ |      $$ \__$$ |$$ |$$$$$$$$/
$$$/    $$$ |$$    $$/ $$ |      $$    $$ |$$ |$$       |
$$/      $$/  $$$$$$/  $$/        $$$$$$$/ $$/  $$$$$$$/ `)

	fmt.Println("")
	fmt.Println(color.BlueString("%s", "- A Blue letter is a letter in the wrong place"))
	fmt.Println(color.GreenString("%s", "- A Green letter is a letter in the right place"))
	fmt.Println("- A White letter is a letter that isn't in the word")
	fmt.Println("")
	fmt.Println("Guess the word: _ _ _ _ _")

	for !game.ended {
		var guess string

		fmt.Print("> ")
		fmt.Scanln(&guess)

		guessArr := strings.Split(guess, "")

		if guess != "" && len(guessArr) == 5 {
			if game.try != 4 {
				clearConsole()
			}
			game.EndRound(guess)
		}
	}

	if game.won {
		fmt.Printf("\n\nYou won! The word: %s is correct", game.word)
	} else if !game.won {
		fmt.Printf("\n\nYou lost. The correct word was: %s.", game.word)
	}
}
