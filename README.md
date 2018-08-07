# Prerequisites
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Go](https://golang.org/doc/install)
- [Python 3](https://www.python.org/downloads/)
- [matplotlib](https://matplotlib.org/users/installing.html)

# Quick Start
```
cd $(go env GOPATH)/src
git clone https://github.com/larkuzo/jwt-compare.git
cd jwt-compare
go run main.go > result.csv
python3 result.py
```
