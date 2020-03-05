slogs
===

```
NAME:
   slogs - to squeeze log files

USAGE:
   slogs [global options] command [command options] [FilePath or DirPath (default: ./)]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ext value   extension (default: ".txt")
   -r            recursion flag (default: false)
   --rm          remove files (default: false)
   --dlen value  date length (default: 15)
   --dpat value  date pattern. See https://golang.org/src/time/format.go (default: "Jan _2 15:04:05")
   --as value    file size (more or equal) for async, MB (default: 5)
   --help, -h    show help (default: false)
```

## Example

Input file **19102019.txt**
```text
Oct 19 01:20:58 8.8.8.8 CRIT: Optical transceiver supply voltage low alarm threshold exceeded
Oct 19 01:21:50 8.8.8.8 WARN: Optical transceiver supply voltage back to normal
Oct 19 14:21:53 8.8.8.8 WARN: Optical transceiver RX power back to normal
Oct 19 14:23:07 8.8.8.8 WARN: Optical transceiver RX power low warning threshold exceeded
Oct 19 14:23:57 8.8.8.8 WARN: Optical transceiver RX power back to normal
Oct 19 14:25:16 8.8.8.8 WARN: Optical transceiver RX power low warning threshold exceeded
Oct 19 18:27:19 8.8.8.8 WARN: Optical transceiver RX power back to normal
Oct 19 18:29:50 8.8.8.8 CRIT: Optical transceiver RX power high alarm threshold exceeded
Oct 19 18:29:50 8.8.8.8 CRIT: Optical transceiver TX power high alarm threshold exceeded
Oct 19 18:29:50 8.8.8.8 CRIT: Optical transceiver supply voltage low alarm threshold exceeded
Oct 19 18:30:42 8.8.8.8 WARN: Optical transceiver RX power back to normal
Oct 19 18:30:42 8.8.8.8 WARN: Optical transceiver TX power back to normal
```

Command
```bash
slogs 19102019.txt
```

Output file **_19102019.txt**
```text
Oct 19 01:21:50 8.8.8.8 WARN: Optical transceiver supply voltage back to normal
Oct 19 14:25:16 8.8.8.8 WARN: Optical transceiver RX power low warning threshold exceeded
Oct 19 18:29:50 8.8.8.8 CRIT: Optical transceiver RX power high alarm threshold exceeded
Oct 19 18:29:50 8.8.8.8 CRIT: Optical transceiver TX power high alarm threshold exceeded
Oct 19 18:29:50 8.8.8.8 CRIT: Optical transceiver supply voltage low alarm threshold exceeded
Oct 19 18:30:42 8.8.8.8 WARN: Optical transceiver RX power back to normal
Oct 19 18:30:42 8.8.8.8 WARN: Optical transceiver TX power back to normal

==================== Lines more than 2 ====================

4	8.8.8.8 WARN: Optical transceiver RX power back to normal
2	8.8.8.8 CRIT: Optical transceiver supply voltage low alarm threshold exceeded
2	8.8.8.8 WARN: Optical transceiver RX power low warning threshold exceeded
```

### Other examples
```bash
slogs -r logs/
```

```bash
slogs -r -rm logs/
```

```bash
slogs -rm logs/19102019.txt
```