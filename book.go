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
	"errors"
	"fmt"
)

// Account stores the information of an account in bookkeeping.
type Account struct {
	Name        string
	EntryFile   string
	liabilities map[string]float64
}

// Book stores the information of a book.
type Book struct {
	IsSigned   bool
	Currencies []string
	Accounts   []Account
}

// GetAccountByName returns the account of name or error occurs
// when the account was not found.
func (b *Book) GetAccountByName(name string) (*Account, error) {
	for _, ac := range b.Accounts {
		if ac.Name == name {
			return &ac, nil
		}
	}

	return nil, fmt.Errorf("requested account %s not found", name)
}

// GetLiabilities returns the liabilities in curs of an account.
// Empty parameter list will return all liabilities.
// Error occurs when any string in curs is not match any valid currencies.
func (ac *Account) GetLiabilities(curs ...string) (map[string]float64, error) {
	if ac.liabilities == nil {
		return nil, errors.New("validation must be excuted before getting liabilities")
	}

	if len(curs) == 0 {
		return ac.liabilities, nil
	}

	l := make(map[string]float64)
	for _, c := range curs {
		if v, ok := ac.liabilities[c]; ok {
			l[c] = v
		} else {
			return nil, fmt.Errorf("currency %s not found", c)
		}
	}

	return l, nil
}

func (b *Book) checkCurrency(s string) bool {
	for _, c := range b.Currencies {
		if c == s {
			return true
		}
	}

	return false
}
