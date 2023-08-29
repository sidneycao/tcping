module github.com/sidneycao/tcping

go 1.18

require github.com/spf13/cobra v1.7.0

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

retract (
    v1.0.0 // fix repo path
    [v0.0.2, v0.0.7] // Retract all from v0
)