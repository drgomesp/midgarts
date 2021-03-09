# Midgarts

Midgarts Client is an attempt to write a modern client implementation of the old classic [Ragnarök Online](https://ragnarok.fandom.com/wiki/Ragnarok_Online) game. 

## Table of Contents

- Introduction (coming soon)
- [Examples](https://github.com/drgomesp/midgarts/blob/master/README.md#examples)

- [TODO](https://github.com/drgomesp/midgarts/blob/master/TODO.md#todo)

## Introduction

## Examples

### Loading a GRF file

```go
grfFilef, err := grf.Load("data.grf")
```

### Getting an entry

```go
grfEntry, err := f.GetEntry("data\sprite\ork_warrior.spr")
```

### Loading SPR files

```go
sprFile, err := spr.Load(e.Data)
```

### Generating a PNG from a sprite

```go
outputFile, err := os.Create("out/test.png")
if err != nil {
    log.Fatal(err)
}

defer outputFile.Close()

if err = png.Encode(outputFile, img); err != nil {
    log.Fatal(err)
}
```
