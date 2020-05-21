Poorloan
========

Poorloan is something to play with friends. We love borrowing money from friends, and forgetting them after spending. But our friends are nice, we should not betray them, and should mark all of your credit on the wall.

All of us are poor, and we love this poor program. The main goal is to have every transaction recorded, with simple currency conversion. 

## Setup

Nothing else but an editor is required for accounting, even without any accounting experience. Poorloan is just a package providing validation on book and computing liabilities.

The book will be separated into files. A `.bk` file is the main file containing book options and all basic information of accounts. A lot of `.ent` files do the accounting, containing only entries related to the particular account.

### Book Options

#### Currency Setting

    currencies {string} ...

`{}` blocks refer to variables here and below, and the variable type is specified; otherwise the word is of keywords. Many `{string}` blocks can be appended behind. `string` accepts no blank characters (spaces or tabs), due to blank characters are separators of blocks.

Only the last one `currencies` command is valid in the book.

#### Account Setting

    account {string} {string}/{quoted string}

Each `account` command opens a new account. The first `string` refers to the unique name of the account, followed by the filename which is the `.ent` file associated with the account. A string quoted with `""` allows the filename containing spaces (maybe buggy if too complicated).

### Entries

Entries are simple and we do double-entry accounting between accounts.

    {uuid} credit/debit {Decimal} {string} [ @/= {Decimal} {string} ]
      ID                  Amount  Currency         Amount  Currency

Any Currency `string` blocks must be a string registered with `currencies` command in `.bk` file. ID in uuid form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx is unique for every transaction. In double-entry accounting, all entries of a single ID must be a **zero** sum when calculating the amount (with sign, debit is positive and credit negative). The first `Amount` is the actual transaction amount, of the first `Currency`. `[]` blocks are optional. Here the optional block is used to do currency conversion while calculating account liabilities.

The concept here is that we do not *returning* the money back, instead we credit for filling debt. One holds negative liabilities **should** have the responsibility to credit money to others. In this accounting program, `debit` means the holder of this account receives money, while `credit` means sending. The one who receives money gets negative liabilities. 

Considering human-readability, we use `credit` or `debit` to record the entry instead of simple plus/minus sign. Therefore, `Amount` must be a positive number, while the program does not check due to flexibility. Not using plus/minus sign also give convenience on writing currency conversion information.

#### Currency Conversion

Currency conversion use two operator: `@` and `=`, referring to "convert at price" and "convert to amount". For example, 

    10 USD @ 7 CNY

tells the program to write 70 CNY to the account liabilities, while

    10 USD = 70 CNY

does the same. It is easy to receive money then, using `@`,  automatically get the liability of another currency with the current exchange rate; as well as fill the debt with specified amount using `=` when *returning*. Just feel free to use any style in your cases.

## The Package

This **Go** package can be installed typing

    go get github.com/GreenYun/poorloan

then get the usage detail typing (after installing)

    go doc -all poorloan

## Future Development

To work with GPG signature is on the development agenda. The `signature on/off` option is preserved for future use. This accounting method is designed for distributed accounting that every single man keeps his account and, if these guys do not trust each other very much, for every entry in the file, the stakeholders use their private key to sign it. 

## If you are interested in ...

It is very easy that you just simply install [Git](https://git-scm.com/) and [Go](https://golang.org/) on your machine following their installation instruction. Then you type in the command line:

    git clone https://github.com/GreenYun/poorloan.git
    cd poorloan
    git checkout development-master

and you are in developer mode. Read `README.md` in the `test` folder to know what to do with the test code. Feel free to make enhancement, new [pull requests](https://github.com/GreenYun/poorloan/pulls) and new [issues](https://github.com/GreenYun/poorloan/issues), even new programs to automatically process bookkeeping.