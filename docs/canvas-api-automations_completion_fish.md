## canvas-api-automations completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	canvas-api-automations completion fish | source

To load completions for every new session, execute once:

	canvas-api-automations completion fish > ~/.config/fish/completions/canvas-api-automations.fish

You will need to start a new shell for this setup to take effect.


```
canvas-api-automations completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --courseID int      Specify course ID, necessary for most sub-commands
  -l, --logLevel string   Sets log level (fatal|error|warn|info|debug|trace) (default "info")
      --readOnly          Set to disable all non-GET http requests such as POST and PUT
```

### SEE ALSO

* [canvas-api-automations completion](canvas-api-automations_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 12-Feb-2024