package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"../../poorloan"
	"github.com/google/uuid"
)

const bookFileName = "test.bk"

func main() {
	fNoMakeData := flag.Bool("no-make-data", false, "Test only, do not make test data.")
	fNoTest := flag.Bool("no-test", false, "Make test data only, do not run the test.")
	flag.Parse()

	var testUUID uuid.UUID
	if !*fNoMakeData {
		testUUID = mktestdata()
	}

	if !*fNoTest {
		test(testUUID)
	}
}

func mktestdata() uuid.UUID {
	fb, err := os.Create(bookFileName)
	if err != nil {
		panic(err)
	}
	currencies := []string{"AUD", "EUR", "JPY", "USD"}
	fb.WriteString("currencies " + strings.Join(currencies, " ") + "\n")

	rand.Seed(time.Now().UnixNano())

	nFile := rand.Intn(6) + 2
	files := make([]*os.File, nFile)
	for i := range files {
		fn := "acc" + strconv.FormatInt(int64(i), 10) + ".ent"
		f, err := os.Create(fn)
		if err != nil {
			panic(err)
		}

		files[i] = f
		fb.WriteString(fmt.Sprintf("account account%d \"%s\"\n", i, fn))
	}

	fb.Sync()
	fb.Close()

	var ret uuid.UUID
	for i := 2; i <= nFile; i++ {
		for j := 0; j < 1000; j++ {
			id, err := uuid.NewRandom()
			if err != nil {
				panic(err)
			}

			if rand.Intn(100) == nFile {
				ret = id
			}

			cur := currencies[rand.Intn(4)]

			amount := make([]float64, i)
			var sum float64 = 0
			for k := 0; k < i-1; k++ {
				if cur == "JPY" {
					amount[k] = (float64(rand.Intn(2)) - 0.5) * float64(rand.Int31()/1000) * 2
				} else {
					amount[k] = (float64(rand.Intn(2)) - 0.5) * float64(rand.Int31()/1000) / float64(50)
				}
				sum += amount[k]
			}
			amount[i-1] = -sum

			selected := make([]bool, nFile)
			for k := 0; k < i; k++ {
				sf := rand.Intn(nFile)
				for ; selected[sf]; sf++ {
					if sf == nFile {
						sf = -1
					}
				}

				var transfer string = "debit"
				if amount[k] < 0 {
					transfer = "credit"
				}

				convert := rand.Intn(3)
				midval := rand.Intn(20) - 10
				if midval < 0 {
					midval = 0
				}
				convertAmount := rand.Float64() + float64(midval)
				convertCur := currencies[rand.Intn(4)]
				if cur == convertCur {
					convert = 0
				}

				switch convert {
				case 0:
					files[sf].WriteString(fmt.Sprintf("%s %s %.2f %s\n", id.String(), transfer, math.Abs(amount[k]), cur))
				case 1:
					if convertCur == "JPY" {
						files[sf].WriteString(fmt.Sprintf("%s %s %.2f %s @ %.2f %s 1\n", id.String(), transfer, math.Abs(amount[k]), cur, convertAmount, convertCur))
					} else {
						files[sf].WriteString(fmt.Sprintf("%s %s %.2f %s @ %.2f %s\n", id.String(), transfer, math.Abs(amount[k]), cur, convertAmount, convertCur))
					}
				case 2:
					if convertCur == "JPY" {
						files[sf].WriteString(fmt.Sprintf("%s %s %.2f %s = %.2f %s\n", id.String(), transfer, math.Abs(amount[k]), cur, math.Round(convertAmount*math.Abs(amount[k])), convertCur))
					} else {
						files[sf].WriteString(fmt.Sprintf("%s %s %.2f %s = %.2f %s\n", id.String(), transfer, math.Abs(amount[k]), cur, convertAmount*math.Abs(amount[k]), convertCur))
					}
				}
			}
		}
	}

	for _, f := range files {
		f.Sync()
		f.Close()
	}

	return ret
}

func test(id uuid.UUID) {
	b, err := poorloan.Parse("test.bk")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Currencies: ", b.Currencies)
	fmt.Println("\nAccount name, filename")
	for _, ac := range b.Accounts {
		fmt.Println(ac.Name, ac.EntryFile)
	}

	ent, err := b.ValidateEntryAndGet(id)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nGet an entry: ", ent)

	fmt.Println("\nPrint all liabilities")
	for _, v := range b.Accounts {
		fmt.Println(v.GetLiabilities())
	}
}
