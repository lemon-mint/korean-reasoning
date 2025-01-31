package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"koreanreasoning/internal/types"
	"os"
	"strings"
)

func main() {
	inputFile := "data/open-thoughts-114k.jsonl"
	outputFile := "data/ot-processed-01.jsonl"
	outputFailFile := "data/ot-processed-01.jsonl.fail"

	readFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer readFile.Close()

	writeFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writeFile.Close()

	writeFailFile, err := os.Create(outputFailFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer writeFailFile.Close()

	reader := bufio.NewReader(readFile)
	writer := bufio.NewWriter(writeFile)
	writerFail := bufio.NewWriter(writeFailFile)
	defer writer.Flush()
	defer writerFail.Flush()

	line_no := 0
	for {
		line_no++
		if line_no%1000 == 0 {
			fmt.Printf("Processing line %d\n", line_no)
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		var ot114kLine types.OT114K_LINE
		err = json.Unmarshal([]byte(line), &ot114kLine)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue
		}

		ok, processedLine := process_line(&ot114kLine)
		if ok {
			processedLineJSON, err := json.Marshal(processedLine)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				continue
			}
			_, err = writer.WriteString(string(processedLineJSON) + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		} else {
			writerFail.WriteString(line)
			writerFail.WriteString("\n")
		}
	}

	fmt.Println("Processing complete. Output saved to", outputFile)
}

var prefixes = []string{
	"Generate an executable Python function generated from the given prompt. The function should take stdin as input and print the output. Simply call the function after the definition.",
	"Return your final response within \\boxed{}. ",
	"Generate an executable Python function generated from the given prompt. Return the function body without invoking it at the final solution.",
}

func process_line(in *types.OT114K_LINE) (ok bool, out *types.OT_PROCESSED_01_LINE) {
	if len(in.Conversations) == 2 {
		// 1. role check
		if in.Conversations[0].From != "user" {
			fmt.Println("Fail 11")
			return false, nil
		}

		if in.Conversations[1].From != "assistant" {
			fmt.Println("Fail 12")
			return false, nil
		}

		question := in.Conversations[0].Value
		generated := in.Conversations[1].Value

		// 2. remove prefix
		removed := false
		prefix_id := -1
		for i := range prefixes {
			if strings.HasPrefix(question, prefixes[i]) {
				question = strings.TrimPrefix(question, prefixes[i])
				prefix_id = i
				removed = true
			}
		}

		// 3. remove quotes (pre)
		if !removed && strings.HasPrefix(question, "\"") {
			if !strings.HasSuffix(question, "\"") {
				fmt.Println("Fail 31")
				return false, nil
			}

			question = strings.TrimPrefix(question, "\"")
			question = strings.TrimSuffix(question, "\"")
		}

		// 4. Clean Whitespace
		question = strings.TrimSpace(question)

		// 5. check tokens
		if !strings.HasPrefix(generated, "<|begin_of_thought|>") {
			fmt.Println("Fail 51")
			return false, nil
		}
		generated = strings.TrimPrefix(generated, "<|begin_of_thought|>")

		if !strings.HasSuffix(generated, "<|end_of_solution|>") {
			fmt.Println("Fail 52")
			return false, nil
		}
		generated = strings.TrimSuffix(generated, "<|end_of_solution|>")

		reasoning, response, ok := strings.Cut(generated, "<|end_of_thought|>")
		if !ok {
			fmt.Println("Fail 53")
			return false, nil
		}

		dummy, response, ok := strings.Cut(response, "<|begin_of_solution|>")
		if !ok {
			fmt.Println("Fail 54")
			return false, nil
		}

		if strings.TrimSpace(dummy) != "" {
			fmt.Println("Fail 55")
			return false, nil
		}

		reasoning = strings.TrimSpace(reasoning)
		response = strings.TrimSpace(response)

		// 6. remove quotes (post)

		if strings.HasPrefix(reasoning, "\"") {
			if !strings.HasSuffix(reasoning, "\"") {
				fmt.Println("Fail 61")
				return false, nil
			}

			reasoning = strings.TrimPrefix(reasoning, "\"")
			reasoning = strings.TrimSuffix(reasoning, "\"")
		}

		if strings.HasPrefix(response, "\"") {
			if !strings.HasSuffix(response, "\"") {
				fmt.Println("Fail 62")
				return false, nil
			}

			response = strings.TrimPrefix(response, "\"")
			response = strings.TrimSuffix(response, "\"")
		}

		return true, &types.OT_PROCESSED_01_LINE{
			PrefixID:  prefix_id,
			Question:  question,
			Reasoning: reasoning,
			Response:  response,
		}
	}
	return false, nil
}
