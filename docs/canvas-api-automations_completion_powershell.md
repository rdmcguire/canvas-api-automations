## canvas-api-automations completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	canvas-api-automations completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
canvas-api-automations completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
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

###### Auto generated by spf13/cobra on 22-Apr-2024
