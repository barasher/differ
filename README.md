# Differ

## Presentation

**Differ** is a tools that compares two folders.
- Used as an API, it returns the differences.
- Used as a standalone program, it prints the differentes.

## Use-case

I have developped **Differ** to check that I don't miss anything when I perform manual backup (on my NAS).
I don't like the principle of synchronizing folders because of the unattented deletion risk.

I usually copy my important files from folder A to folder B and I check with **Differ** if something new has not
beenbackuped yet or if something has changed an not been backuped yet

## Execution

```
./differ [-c {configurationFile}] {from} {to}
```

* `-c {configurationFile}` : configuration file
* `{from}` : 'source' folder
* `{to}` : 'target' folder

Return codes :
* 0 : no difference
* 1 : differences detected
* 2 : configuration error
* 3 : execution error

## Configuration file

If no configuration file is provided :
* there will be no blacklisted files

Sample configuration file :
```json
{
	"BlacklistedPatterns":[
		"^.*txt$",
		"^.*doc$"
	]
}
```

Configuration settings :
* **BlacklistedPatterns** : regular expressions to blacklist files. Blacklisted files won't be checked.

## Output

```
WARN[0000] [S] folder/different.txt
WARN[0000] [M] folder/nonExistingInTo.txt
WARN[0000] [M] folderNonExistingInTo
WARN[0000] [S] different.txt
WARN[0000] [T] differentType
WARN[0000] [M] nonExistingInTo.txt
```

* `[S]` : element has a different size between source and target folder
* `[M]` : element is missing in the target folder
* `[T]` : element has a different type (file / folder) between source and target folder