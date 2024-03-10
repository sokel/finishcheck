package main

import (
	"github.com/sokel/finishcheck/pkg"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(pkg.NewAnalyzer())
}
