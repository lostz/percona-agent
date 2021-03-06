/*
   Copyright (c) 2014, Percona LLC and/or its affiliates. All rights reserved.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

type InvalidResponseError struct {
	Response string
}

func (e InvalidResponseError) Error() string {
	return e.Response
}

type Terminal struct {
	stdin *bufio.Reader
	flags Flags
}

func NewTerminal(stdin io.Reader, flags Flags) *Terminal {
	t := &Terminal{
		stdin: bufio.NewReader(stdin),
		flags: flags,
	}
	return t
}

func (t *Terminal) PromptString(question string, defaultAnswer string) (string, error) {
	if defaultAnswer != "" {
		fmt.Printf("%s (%s): ", question, defaultAnswer)
	} else {
		fmt.Printf("%s: ", question)
	}
	bytes, _, err := t.stdin.ReadLine()
	if err != nil {
		return "", err
	}
	if t.flags["debug"] {
		log.Printf("raw answer='%s'\n", string(bytes))
	}
	answer := strings.TrimSpace(string(bytes))
	if answer == "" {
		answer = defaultAnswer
	}
	if t.flags["debug"] {
		log.Printf("final answer='%s'\n", answer)
	}
	return answer, nil
}

func (t *Terminal) PromptStringRequired(question string, defaultAnswer string) (string, error) {
	var answer string
	var err error
	for {
		answer, err = t.PromptString(question, defaultAnswer)
		if err != nil {
			return "", err
		}
		if answer == "" {
			fmt.Println(question + " is required, please try again")
			continue
		}
		return answer, nil
	}
}

func (t *Terminal) PromptBool(question string, defaultAnswer string) (bool, error) {
	for {
		answer, err := t.PromptString(question, defaultAnswer)
		if t.flags["debug"] {
			log.Printf("again=%t\n", answer)
			log.Printf("err=%s\n", err)
		}
		if err != nil {
			return false, err
		}
		answer = strings.ToLower(answer)
		switch answer {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			log.Println("Invalid response: '" + answer + "'.  Enter 'y' for yes, 'n' for no.")
			continue
		}
	}
}
