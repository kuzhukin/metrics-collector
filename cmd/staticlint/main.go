package main

/*
# Использование:

	staticlint [package]

# Проверить все файлы во всех директориях проекта:

	staticlint ./...

# Справка

	staticlint help

# Пакет представляет собой набор следующих статических анализаторов (по умолчанию - без файла конфигурации):
  - стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes
  - анализаторы класса SA пакета staticcheck:
    SA1* - проблемы исопльзования стандартной библиотека,
    SA2* - проблемы с конкурентностью,
    SA3* - проблемы с тестированием,
    SA4* - бесполезный код (который ничего не делает по факту)
    SA5* - проблемы с корректностью
    SA6* - Проблемы с производительностью
    SA9* - Сомнительные кодовые конструкции (с высокой вероятностью могут быть ошибочными)
  - анализаторы класса S пакета staticcheck:
    S1001 - замена цикла for на вызов copy() для слайсов
  - анализаторы класса ST пакета staticcheck:
    ST1000 - неправильный или отсутствующий комментарий к пакету
  - анализаторы класса QF пакета staticcheck:
    QF1003 - преобразование if/else-if в switch
  - анализатор `github.com/Antonboom/errname`:
    проверяет, что переменные ошибок имеют префикс Err, а типы ошибок - суффикс Error
  - анализатор `github.com/leonklingele/grouper`:
    анализ групп выражений (импорты, типы, переменные, константы)
  - анализатор `github.com/sashamelentyev/usestdlibvars`:
    обнаруживает возможность использования переменных/констант из стандартной библиотеки
  - анализатор noexitanalyzer, проверяет использование прямого вызова os.Exit в функции main
*/

import (
	"github.com/kuzhukin/metrics-collector/cmd/staticlint/noexitanalyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

var saChecks = []string{
	"S1001",
	"S1002",
	"S1025",
	"S1028",
	"ST1000",
	"ST1005",
	"ST1006",
	"ST1017",
	"ST1020",
	"ST1021",
	"ST1022",
	"ST1023",
	"QF1003",
	"QF1006",
	"QF1007",
}

func main() {
	linters := []*analysis.Analyzer{
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpmux.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
		// custom analyzers
		noexitanalyzer.Analyzer,
	}

	staticchecks := make(map[string]bool)

	for _, v := range saChecks {
		staticchecks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		linters = append(linters, v.Analyzer)
	}

	for _, v := range quickfix.Analyzers {
		if staticchecks[v.Analyzer.Name] {
			linters = append(linters, v.Analyzer)
		}
	}

	for _, v := range simple.Analyzers {
		if staticchecks[v.Analyzer.Name] {
			linters = append(linters, v.Analyzer)
		}
	}

	for _, v := range stylecheck.Analyzers {
		if staticchecks[v.Analyzer.Name] {
			linters = append(linters, v.Analyzer)
		}
	}

	multichecker.Main(linters...)
}
