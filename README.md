# `tint`: 🌈 **slog.Handler** that writes tinted logs

[![Go Reference](https://pkg.go.dev/badge/github.com/phplego/tint.svg)](https://pkg.go.dev/github.com/phplego/tint#section-documentation)
[![Go Report Card](https://goreportcard.com/badge/github.com/phplego/tint)](https://goreportcard.com/report/github.com/phplego/tint)


![image](https://github.com/user-attachments/assets/ad3855ea-ab17-4109-b027-3f0fbb67e0a2)

<br>
<br>

Package `tint` implements a zero-dependency [`slog.Handler`](https://pkg.go.dev/log/slog#Handler)
that writes tinted (colorized) logs. Its output format is inspired by the `zerolog.ConsoleWriter` and
[`slog.TextHandler`](https://pkg.go.dev/log/slog#TextHandler).

The output format can be customized using [`Options`](https://pkg.go.dev/github.com/lmittmann/tint#Options)
which is a drop-in replacement for [`slog.HandlerOptions`](https://pkg.go.dev/log/slog#HandlerOptions).

```
go get github.com/phplego/tint
```

## Usage

```go
w := os.Stderr

// create a new logger
logger := slog.New(tint.NewHandler(w, nil))

// set global logger with custom options
slog.SetDefault(slog.New(
    tint.NewHandler(w, &tint.Options{
        Level:      slog.LevelDebug,
        TimeFormat: time.Kitchen,
    }),
))

// use color expressions
slog.Info("Message with @r{red color}")
slog.Info("Message with @R{bright red color}")
slog.Info("Message with @r!{red bold color}")
slog.Info("Message with @rW{red on white backgound}")
```

### Customize Attributes

`ReplaceAttr` can be used to alter or drop attributes. If set, it is called on
each non-group attribute before it is logged. See [`slog.HandlerOptions`](https://pkg.go.dev/log/slog#HandlerOptions)
for details.

```go
// create a new logger that doesn't write the time
w := os.Stderr
logger := slog.New(
    tint.NewHandler(w, &tint.Options{
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if a.Key == slog.TimeKey && len(groups) == 0 {
                return slog.Attr{}
            }
            return a
        },
    }),
)
```

```go
// create a new logger that writes all errors in red
w := os.Stderr
logger := slog.New(
    tint.NewHandler(w, &tint.Options{
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if err, ok := a.Value.Any().(error); ok {
                aErr := tint.Err(err)
                aErr.Key = a.Key
                return aErr
            }
            return a
        },
    }),
)
```

### Automatically Enable Colors

Colors are enabled by default and can be disabled using the `Options.NoColor`
attribute. To automatically enable colors based on the terminal capabilities,
use e.g. the [`go-isatty`](https://github.com/mattn/go-isatty) package.

```go
w := os.Stderr
logger := slog.New(
    tint.NewHandler(w, &tint.Options{
        NoColor: !isatty.IsTerminal(w.Fd()),
    }),
)
```

### Windows Support

Color support on Windows can be added by using e.g. the
[`go-colorable`](https://github.com/mattn/go-colorable) package.

```go
w := os.Stderr
logger := slog.New(
    tint.NewHandler(colorable.NewColorable(w), nil),
)
```
