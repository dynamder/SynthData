# Quick Start: Interactive CLI Mode

## Basic Usage

Run any command with the `-i` flag to enter interactive mode:

```bash
synthdata generate -i
```

## Step-by-Step Example

```bash
# Start interactive mode
$ synthdata generate -i

Welcome to Interactive Mode!
Press Ctrl+C at any time to exit.

▶ Description file (-d, --description):
  > schema.json
  ✓ Valid path

▶ Output file (-o, --output):
  > output.json

▶ Scale (-s, --scale) [default: 10]:
  > 100

▶ Format (-f, --format) [json, csv]:
  > json

Collected arguments:
  --description schema.json
  --output output.json
  --scale 100
  --format json

Generating 100 records...
✓ Successfully generated 100 records to output.json
```

## Combining with Explicit Arguments

Skip prompts for args you already know:

```bash
synthdata generate -i -o output.json -f csv
# Only prompts for description file and scale
```

## Getting Help

During any prompt, type `?` or `help` to see detailed information about that argument.

## Exiting

Press `Ctrl+C` at any time to exit. No changes will be made.
