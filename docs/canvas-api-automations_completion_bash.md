## canvas-api-automations completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(canvas-api-automations completion bash)

To load completions for every new session, execute once:

#### Linux:

	canvas-api-automations completion bash > /etc/bash_completion.d/canvas-api-automations

#### macOS:

	canvas-api-automations completion bash > $(brew --prefix)/etc/bash_completion.d/canvas-api-automations

You will need to start a new shell for this setup to take effect.


```
canvas-api-automations completion bash
```

### Options

```
  -h, --help              help for bash
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