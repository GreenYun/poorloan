Test program for poorloan
=========================

The `main.go` is a simple program to test if `poorloan` works fine. To run the program, type

    go run ./main.go

and check the output if there are errors. The program imports `poorloan` in relative path that always runs the latest code, and makes it easy to set breakpoints to debug.

Two functions `mktestdata` and `test` run in turn. `mktestdata` generates random book data for testing. `test` to validate the data and print some patterned output or error. 