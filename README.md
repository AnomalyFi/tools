# Tools for running and testing SEQ.

## Relayer tools:

Creates config file for testing decRen in devnet.
There will be a better way for doing this, but this is the easiest find.
note: Run only SEQ on the machine. not any other chain.

usage :
```GO
go run main.go config.json
```
## Oracle tools:

Creates config file for oracle.
note: Atleast one SEQ node should be run on the machine using this.

usage:
```GO
go run main.go config.json
```