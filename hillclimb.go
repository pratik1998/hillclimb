package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/mkideal/cli"
	"github.com/pratik1998/hillclimb/enigma"
	"github.com/pratik1998/hillclimb/scores"
)

type stringDecoder struct {
	list []string
}

func (d *stringDecoder) Decode(s string) error {
	d.list = strings.Split(s, ",")
	return nil
}

// type CLIOpts struct {
// 	Help      bool `cli:"!h,help" usage:"Show help."`
// 	Condensed bool `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`

// 	Rotors    stringDecoder `cli:"rotors" name:"I II III" usage:"Rotor configuration. Supported: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma."`
// 	Rings     stringDecoder `cli:"rings" name:"1 1 1" usage:"Rotor rings offset: from 1 (default) to 26 for each rotor."`
// 	Position  stringDecoder `cli:"position" name:"A A A" usage:"Starting position of the rotors: from A (default) to Z for each."`
// 	Plugboard stringDecoder `cli:"plugboard" name:"[]" usage:"Optional plugboard pairs to scramble the message further."`

// 	Reflector string `cli:"reflector" name:"C" usage:"Reflector. Supported: A, B, C, B-Thin, C-Thin."`
// }

type HillClimbOpts struct {
	Help bool `cli:"!h,help" usage:"Show help."`
}

type RotorSettings struct {
	leftmost_rotor           string
	second_leftmost_rotor    string
	leftmost_position        int
	second_leftmost_position int
	score                    float64
}

func NewRotorSettings(a string, b string, c int, d int, e float64) *RotorSettings {
	rs := &RotorSettings{}
	rs.leftmost_rotor = a
	rs.second_leftmost_rotor = b
	rs.leftmost_position = c
	rs.second_leftmost_position = d
	rs.score = e
	return rs
}

func (rs *RotorSettings) Print() {
	fmt.Println(rs.leftmost_rotor, rs.second_leftmost_rotor)
	fmt.Println(rs.leftmost_position, rs.second_leftmost_position)
	fmt.Println(rs.score)
}

func dequeue(queue []RotorSettings) []RotorSettings {
	if len(queue) < 0 {
		return queue
	}
	return queue[1:]
}

func enqueue(queue []RotorSettings, element RotorSettings) []RotorSettings {
	if len(queue) >= 5 { // Number of rotors to run hillclimb
		queue = dequeue(queue)
	}
	return append(queue, element)
}

func main() {
	// ciphertext := "PQSPPKXGVGEVBXLLDHTFXJJHNUZHNHCIQQJABCF"
	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(HillClimbOpts), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*HillClimbOpts)
		filename := strings.Join(ctx.Args(), " ")
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading a file: " + filename)
		}
		originalPlaintext := string(b)
		// fmt.Println("Original Plain Text: " + originalPlaintext)
		plaintext := enigma.SanitizePlaintext(originalPlaintext)
		// fmt.Println("Sanitized Plain Text: " + plaintext)

		if argv.Help || len(plaintext) == 0 {
			com := ctx.Command()
			// com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}

		// Trigram scorer
		trigram_scorer := scores.NewTrigramScore()

		// Set up rotors
		rotor_config := make([]enigma.RotorConfig, 4)
		// these two will not be overwritten
		rotor_config[3] = enigma.RotorConfig{"III", 'Q', 16}
		rotor_config[2] = enigma.RotorConfig{"IV", 'B', 1}

		var plugboard []string
		var queue []RotorSettings

		available_rotors := [6]string{"I", "II", "V", "VI", "Beta", "Gamma"}
		var max_ioc_score float64
		for _, leftmost_rotor := range available_rotors {
			if leftmost_rotor != "Beta" && leftmost_rotor != "Gamma" {
				continue
			} else {
				for _, second_leftmost_rotor := range available_rotors {
					if leftmost_rotor == second_leftmost_rotor {
						continue
					}
					for leftmost_position := 0; leftmost_position < 26; leftmost_position++ {
						for second_leftmost_position := 0; second_leftmost_position < 26; second_leftmost_position++ {
							rotor_config[0] = enigma.RotorConfig{leftmost_rotor, enigma.IndexToChar(leftmost_position), 1}
							rotor_config[1] = enigma.RotorConfig{second_leftmost_rotor, enigma.IndexToChar(second_leftmost_position), 1}
							e := enigma.NewEnigma(rotor_config, "C-thin", plugboard)
							encoded := e.EncodeString(plaintext)
							ioc_score := scores.GetIndexOfCoincidence(encoded)
							if max_ioc_score < ioc_score {
								queue = enqueue(queue, *NewRotorSettings(leftmost_rotor, second_leftmost_rotor, leftmost_position, second_leftmost_position, ioc_score))
								max_ioc_score = ioc_score
							}
						}
					}
				}
			}
		}

		// fmt.Println("Best", len(queue), " Rotor Settings:")
		// for _, ele := range queue {
		// 	fmt.Println(ele)
		// }

		var best_rs RotorSettings
		var best_plugs []string
		var best_score float64 = -1000000000

		// start from best rotor settings to worst
		for i := len(queue) - 1; i >= 0; i-- {
			rs := queue[i]
			rotor_config[0] = enigma.RotorConfig{rs.leftmost_rotor, enigma.IndexToChar(rs.leftmost_position), 1}
			rotor_config[1] = enigma.RotorConfig{rs.second_leftmost_rotor, enigma.IndexToChar(rs.second_leftmost_position), 1}
			var plugboard []string
			e := enigma.NewEnigma(rotor_config, "C-thin", plugboard)
			encoded := e.EncodeString(plaintext)
			rs.score = trigram_scorer.GetTrigramScore(encoded)
			var plugged [26]bool
			// Hillclimbing
			// Need to find ten plugboards
			for round := 0; round < 10; round++ {
				// pick two plugs and check statistical property on it
				var bestplug string
				before_score := rs.score
				for j := 0; j < 26; j++ {
					for k := j + 1; k < 26; k++ {
						if plugged[j] || plugged[k] {
							continue
						} else {
							tmp_plugboard := plugboard
							new_plug := string(enigma.IndexToChar(j)) + string(enigma.IndexToChar(k))
							tmp_plugboard = append(tmp_plugboard, new_plug)
							e := enigma.NewEnigma(rotor_config, "C-thin", tmp_plugboard)
							encoded := e.EncodeString(plaintext)
							curr_score := trigram_scorer.GetTrigramScore(encoded)
							if rs.score < curr_score {
								rs.score = curr_score
								bestplug = new_plug
							}
						}
					}
				}
				// No new plugs found
				if before_score == rs.score {
					break
				} else {
					plugged[enigma.CharToIndex(bestplug[0])] = true
					plugged[enigma.CharToIndex(bestplug[1])] = true
					plugboard = append(plugboard, bestplug)
					// fmt.Println("Best Plug after round #", round, " ", bestplug)
				}
			}
			sort.Strings(plugboard)
			// fmt.Println(rs)
			// fmt.Println(plugboard)
			// rotor_config[0] = enigma.RotorConfig{rs.leftmost_rotor, enigma.IndexToChar(rs.leftmost_position), 1}
			// rotor_config[1] = enigma.RotorConfig{rs.second_leftmost_rotor, enigma.IndexToChar(rs.second_leftmost_position), 1}
			// fmt.Println(rotor_config)
			// e = enigma.NewEnigma(rotor_config, "C-thin", plugboard)
			// fmt.Println("Decoded Message: ", e.EncodeString(plaintext))
			if best_score < rs.score {
				best_rs = rs
				best_plugs = plugboard
				best_score = rs.score
			}
		}

		fmt.Print(best_rs.leftmost_rotor + " " + best_rs.second_leftmost_rotor + " IV III\n")
		fmt.Print(string(enigma.IndexToChar(best_rs.leftmost_position)) + " " + string(enigma.IndexToChar(best_rs.second_leftmost_position)) + " B Q\n")
		plugboard_str := ""
		for i, str := range best_plugs {
			if i == len(best_plugs)-1 {
				plugboard_str += (str)
			} else {
				plugboard_str += (str + " ")
			}
		}
		if len(plugboard_str) > 0 {
			fmt.Print(plugboard_str + "\n")
		}

		// config := make([]enigma.RotorConfig, len(argv.Rotors.list))
		// for index, rotor := range argv.Rotors.list {
		// 	ring := argv.Rings.list[index]
		// 	value := argv.Position.list[index][0]
		// 	ringValue, _ := strconv.Atoi(ring)
		// 	config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ringValue}
		// }

		// fmt.Println("Received Enigma Settings")
		// fmt.Println("Rotors:")
		// for index, rotor := range argv.Rotors.list {
		// 	fmt.Println(string(argv.Rings.list[index]) + " " + string(argv.Position.list[index][0]) + " " + rotor)
		// }
		// fmt.Println("Reflector: " + argv.Reflector)
		// fmt.Println("Plugboard:")
		// for index, plugsetting := range argv.Plugboard.list {
		// 	fmt.Println(index, plugsetting)
		// }
		// trigram_scorer := scores.NewTrigramScore()
		// e := enigma.NewEnigma(config, argv.Reflector, argv.Plugboard.list)
		// encoded := e.EncodeString(plaintext)
		// fmt.Println("Sanitized Plain Text: " + plaintext)
		// fmt.Println("Index of Co-incidence for sanitized plain text: ", scores.GetIndexOfCoincidence(plaintext))
		// fmt.Println("Trigram Score for sanitized plain text: ", trigram_scorer.GetTrigramScore(plaintext))
		// fmt.Println("Encoded Message: " + encoded)
		// fmt.Println("Index of Co-incidence for encoded text: ", scores.GetIndexOfCoincidence(encoded))
		// fmt.Println("Trigram Score for encoded text: ", trigram_scorer.GetTrigramScore(encoded))
		return nil
	})
}

/*
  Given Enigma Settings:
    Rotor Setup: ?? ?? IV III
    Initial Start Positions: ?? ?? B Q
    Ringstellung: 1 1 1 16
    Plugboard: ??
    Reflector: C-Thin
*/
