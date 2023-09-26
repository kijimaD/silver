# silver

silver is a library for reproducing the system!

## install

```shell
$ go get github.com/kijimad/silver
```

## example

```go
	tasks := []silver.Task{
		dummy(),
	}
	job := silver.NewJob(tasks)
	job.Run()
}

func dummy() silver.Task {
	t := silver.NewTask("dummy")
	t.SetFuncs(silver.ExecFuncParam{
		TargetCmd: nil,
		DepCmd:    nil,
		InstCmd:   func() error { return t.Exec("echo hello && sleep 2 && echo hello && echo hello") },
	})

	return t
}
```

result

```
[1/1 dummy]
  => [exec] echo hello && sleep 2 && echo hello && echo hello
  => hello
  => hello
  => hello
  => [result] Success install
[1/1 dummy] Success install
```
