# Midgarts

Midgarts Client is an attempt to write a modern client implementation of the old classic [Ragnarök Online](https://ragnarok.fandom.com/wiki/Ragnarok_Online) game. 

Current Screenshots:

![wip](https://user-images.githubusercontent.com/696982/111909241-4a594500-8a3b-11eb-911a-0127da565c2b.png)

## Table of Contents

- Introduction (coming soon)
- [TODO](https://github.com/drgomesp/midgarts/blob/master/TODO.md#todo)
- [Tools](https://github.com/drgomesp/midgarts/blob/master/README.md#tools)
    - [GRF Explorer](https://github.com/drgomesp/midgarts/blob/master/README.md#grf-explorer)
- [Examples](https://github.com/drgomesp/midgarts/blob/master/README.md#examples)

## Introduction

## Tools

### GRF Explorer

Latest screenshots:

![image](https://user-images.githubusercontent.com/696982/111029961-72fb9200-83de-11eb-8707-ded945850305.png)
![image](https://user-images.githubusercontent.com/696982/111030058-0339d700-83df-11eb-8546-0cc931ce36ed.png)


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

