# Midgarts

Midgarts Client is an attempt to write a modern client implementation of the old classic [Ragnarök Online](https://ragnarok.fandom.com/wiki/Ragnarok_Online) game.

Current Screenshots:

<p align="center"">
<img src="https://user-images.githubusercontent.com/696982/115995661-89d2fe00-a5b2-11eb-8801-2eef65d31881.gif" width="44%" />
<img src="https://user-images.githubusercontent.com/696982/115995910-96a42180-a5b3-11eb-8200-1cfae06bf5bc.gif" width="44%" />

## Table of Contents

- Introduction (coming soon)
- [TODO](https://github.com/drgomesp/midgarts/blob/master/TODO.md#todo)
- [Dependencies](https://github.com/drgomesp/midgarts/blob/master/README.md#dependencies)
- [Tools](https://github.com/drgomesp/midgarts/blob/master/README.md#tools)
    - [GRF Explorer](https://github.com/drgomesp/midgarts/blob/master/README.md#grf-explorer)
- [Examples](https://github.com/drgomesp/midgarts/blob/master/README.md#examples)

## Introduction

## Dependencies

1. CentOS/Fedora-like Linux Distros:
   `SDL2{,_image,_mixer,_ttf,_gfx}-devel alsa-lib-devel libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel`

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
