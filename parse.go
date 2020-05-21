/*
	MIT License

	Copyright (C) 2020 GreenYun Organization

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

package poorloan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Parse decodes the fp into a Book or returns an error.
func Parse(fp string) (*Book, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	book := new(Book)
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" || text[0] == '#' {
			continue
		}
		i++
		err = book.ParseLine(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("line %d in %s: %s", i, fp, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return book, nil
}

// ParseLine gets information from l and writes to b or returns an error.
func (b *Book) ParseLine(l string) error {
	s := strings.Split(l, " ")
	s = removeEmpty(&s)

	switch s[0] {
	case "signature":
		if len(s) != 2 {
			return makeSyntaxError(s[0])
		}
		switch s[1] {
		case "on":
			b.IsSigned = true
		case "off":
			b.IsSigned = false
		default:
			return makeSyntaxError(s[0])
		}
	case "currencies":
		if len(s) < 2 {
			return makeSyntaxError(s[0])
		}
		b.Currencies = s[1:]
	case "account":
		splitByQuote := strings.Split(l, "\"")
		if len(splitByQuote) != 3 {
			return makeSyntaxError(s[0])
		}
		if st := strings.TrimSpace(splitByQuote[2]); st != "" {
			return makeSyntaxError(s[0])
		}

		s := strings.Split(splitByQuote[0], " ")
		s = removeEmpty(&s)
		if len(s) != 2 {
			return makeSyntaxError(s[0])
		}

		ac := new(Account)
		ac.Name = s[1]
		ac.EntryFile = splitByQuote[1]
		ac.liabilities = make(map[string]float64)
		b.Accounts = append(b.Accounts, *ac)
	default:
		return errors.New("syntax error")
	}

	return nil
}

func removeEmpty(s *[]string) []string {
	i := 0
	p := *s
	for _, v := range p {
		if st := strings.TrimSpace(v); st != "" {
			p[i] = st
			i++
		}
	}

	return p[0:i]
}

func makeSyntaxError(s string) error {
	return fmt.Errorf("syntax error after ``%s''", s)
}
