package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	FileName string
	Limit    int
	Follow   bool
}

func main() {
	config, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println("erreur:", err)
		return
	}

	lines, err := tailFile(config.FileName, config.Limit)
	if err != nil {
		fmt.Println("erreur:", err)
		return
	}

	for _, line := range lines {
		fmt.Println(line)
	}

	if config.Follow {
		err := followFile(config.FileName)
		if err != nil {
			fmt.Println("erreur:", err)
			return
		}
	}
}

func parseArgs(args []string) (Config, error) {
	config := Config{Limit: 10}

	if len(args) == 0 {
		return config, fmt.Errorf("veuillez renseigner le nom du fichier à analyser")
	}

	for i := 0; i < len(args); {
		arg := args[i]

		switch arg {
		case "-n":
			if i+1 >= len(args) {
				return config, fmt.Errorf("veuillez renseigner le nombre de lignes limite")
			}

			limit, err := strconv.Atoi(args[i+1])
			if err != nil {
				return config, fmt.Errorf("veuillez entrer une valeur numérique")
			}

			if limit <= 0 {
				return config, fmt.Errorf("veuillez entrer une limite entière strictement positive")
			}

			config.Limit = limit
			i += 2
		case "-f":
			config.Follow = true
			i++
		default:
			if strings.HasPrefix(arg, "-") {
				return config, fmt.Errorf("option %s invalide, options disponibles: -n, -f", arg)
			}

			if config.FileName != "" {
				return config, fmt.Errorf("un seul fichier peut être analysé à la fois")
			}

			config.FileName = arg
			i++
		}
	}

	if config.FileName == "" {
		return config, fmt.Errorf("veuillez renseigner le nom du fichier à analyser")
	}

	return config, nil
}

func tailFile(fileName string, limit int) ([]string, error) {
	lastLines := []string{}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lastLines = append(lastLines, line)

		if len(lastLines) > limit {
			lastLines = lastLines[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lastLines, nil
}

func followFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err == nil {
			fmt.Print(line)
			continue
		}

		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		return err
	}
}
