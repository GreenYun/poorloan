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
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Entry stores the information of a single entry.
type Entry struct {
	Credit map[string]float64
	Debit  map[string]float64
}

type bookLine struct {
	index        uuid.UUID
	transferType int
	currency     string
	amount       float64
	asCurrency   string
	asAmount     float64
}

type mutexMap struct {
	m  map[uuid.UUID]float64
	mu sync.Mutex
}

func (mm *mutexMap) plusAssign(key uuid.UUID, val float64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if v, ok := mm.m[key]; ok {
		if sum := v + val; math.Round(sum*100) == 0 {
			delete(mm.m, key)
		} else {
			mm.m[key] = sum
		}
	} else {
		mm.m[key] = val
	}
}

// ValidateEntry validates every entries in the book,
// does the same as ValidateEntryAndGet(uuid.Nil) but returns no entries.
func (b *Book) ValidateEntry() error {
	_, err := b.ValidateEntryAndGet(uuid.Nil)
	return err
}

// ValidateEntryAndGet validates entries in the book,
// and returns the entry of index.
// If index = uuid.Nil, returns nil.
func (b *Book) ValidateEntryAndGet(index uuid.UUID) (*Entry, error) {
	var entryStack mutexMap
	entryStack.m = make(map[uuid.UUID]float64)

	var ret *Entry = nil
	var retMutex sync.Mutex

	g, _ := errgroup.WithContext(context.Background())

	for _, ac := range b.Accounts {
		ac := ac

		g.Go(func() error {
			file, err := os.Open(ac.EntryFile)
			if err != nil {
				return err
			}
			defer file.Close()

			i := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				text := strings.TrimSpace(scanner.Text())
				if text == "" || text[0] == '#' {
					continue
				}

				i++
				bl, err := parseBookLine(text)
				if err != nil {
					return fmt.Errorf("line %d in %s: %s", i, ac.EntryFile, err)
				}
				if !b.checkCurrency(bl.currency) {
					return fmt.Errorf("line %d in %s: currency %s not valid in this book", i, ac.EntryFile, bl.currency)
				}

				entryStack.plusAssign(bl.index, float64(bl.transferType)*bl.amount)

				if bl.asCurrency != "" {
					ac.liabilities[bl.asCurrency] -= float64(bl.transferType) * bl.asAmount
				} else {
					ac.liabilities[bl.currency] -= float64(bl.transferType) * bl.amount
				}

				if bl.index == index {
					retMutex.Lock()

					if ret == nil {
						ret = new(Entry)
						ret.Credit = make(map[string]float64)
						ret.Debit = make(map[string]float64)
					}

					if bl.transferType == 1 {
						ret.Debit[ac.Name] = bl.amount
					} else {
						ret.Credit[ac.Name] = bl.amount
					}

					retMutex.Unlock()
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for k, v := range entryStack.m {
		return nil, fmt.Errorf("entry %s does not meet zero: %f", k.String(), v)
	}

	if index != uuid.Nil && ret == nil {
		return nil, fmt.Errorf("entry %s not found", index.String())
	}

	return ret, nil
}

func parseBookLine(l string) (*bookLine, error) {
	s := strings.Split(l, " ")
	s = removeEmpty(&s)
	translate := false
	if len(s) != 4 {
		if len(s) != 7 {
			return nil, errors.New("syntax error")
		}
		translate = true
	}

	bl := new(bookLine)

	u, err := uuid.Parse(s[0])
	if err != nil {
		return nil, err
	}
	bl.index = u

	switch s[1] {
	case "credit":
		bl.transferType = -1
	case "debit":
		bl.transferType = 1
	default:
		return nil, errors.New("transfer type invalid")
	}

	f64, err := strconv.ParseFloat(s[2], 64)
	if err != nil {
		return nil, err
	}
	bl.amount = f64

	bl.currency = s[3]

	bl.asCurrency = ""
	if translate {
		f64, err = strconv.ParseFloat(s[5], 64)
		if err != nil {
			return nil, err
		}

		switch s[4] {
		case "@":
			bl.asAmount = bl.amount * f64
		case "=":
			bl.asAmount = f64
		default:
			return nil, errors.New("syntax error")
		}

		bl.asCurrency = s[6]
	}

	return bl, nil
}
