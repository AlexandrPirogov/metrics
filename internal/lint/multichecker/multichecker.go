// Package multichecker contains bunch of static lint default checkers
// for SA and ST rules with custom chec
//
// To run multichecker build the source code and specify the dir/file you want to check :)
package multichecker

import (
	"go/ast"
	"log"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
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
	"golang.org/x/tools/go/analysis/passes/slog"
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
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Analyzers that used inside multichecker
// To see their purpose see: golang.org/x/tools/go/analysis :)
var DefaultRules = []*analysis.Analyzer{

	asmdecl.Analyzer,
	assign.Analyzer,
	atomic.Analyzer,
	atomicalign.Analyzer,
	bools.Analyzer,
	buildssa.Analyzer,
	buildtag.Analyzer,
	cgocall.Analyzer,
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
	slog.Analyzer,
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
}

// Regexp for all SA rules
const SAREGEXP = "^[SA]+[\\d]*"

// Regexp for all ST rules
const STREGEXP = "^[ST]+[\\d]*"

// OsExitCheckAnalyzer instance of analysis.Analyzer that checks
// usage of os.Exit inside main function main's package
var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "errexit",
	Doc:  "check for that user does'nt use os.exit in main",
	Run:  run,
}

func New() []*analysis.Analyzer {
	SACompiled, _ := regexp.Compile(SAREGEXP)

	// Appneding SA Rules
	var StaticCheckRules []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if SACompiled.Match([]byte(v.Analyzer.Name)) {
			StaticCheckRules = append(StaticCheckRules, v.Analyzer)
			//log.Printf("Rule %s was added\n", v.Analyzer.Name)
		}
	}

	STCompiler, _ := regexp.Compile(STREGEXP)
	for _, v := range stylecheck.Analyzers {
		if STCompiler.Match([]byte(v.Analyzer.Name)) {
			StaticCheckRules = append(StaticCheckRules, v.Analyzer)
			//	log.Printf("Rule %s was added\n", v.Analyzer.Name)
		}
	}

	// Running golang.org/x/tools/go/analysis/passes default packages

	for _, v := range DefaultRules {
		StaticCheckRules = append(StaticCheckRules, v)
	}

	StaticCheckRules = append(StaticCheckRules, OsExitCheckAnalyzer)
	log.Printf("Static rules: %d\n", len(StaticCheckRules))
	return StaticCheckRules
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if c, ok := n.(*ast.FuncDecl); ok {
				if c.Name.String() == "main" {
					ast.Inspect(c, func(n1 ast.Node) bool {
						if i, ok := n1.(*ast.CallExpr); ok {
							if m, ok := i.Fun.(*ast.SelectorExpr); ok {
								if m.Sel.Name == "Exit" {
									pass.Reportf(i.Pos(), "os.Exit used inside main()")
									return false
								}
							}
						}
						return true
					})
				}
			}

			return true
		})
	}
	return nil, nil
}
