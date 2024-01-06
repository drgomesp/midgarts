# Midgarts

Midgarts Client is an attempt to write a modern client implementation of the old classic [Ragnarök Online](https://ragnarok.fandom.com/wiki/Ragnarok_Online) game.

Current Screenshots:

<p align="center">
    <img src="https://user-images.githubusercontent.com/696982/117575166-fff95980-b0b6-11eb-8afa-acd7dcdd6b34.gif" width="25%" />
    <img src="https://user-images.githubusercontent.com/696982/197043590-041d711b-a5d6-4d58-bf3c-8ea98c1afdc6.gif" width="50%" />
</p>
<p align="center">
    <img src="https://user-images.githubusercontent.com/696982/116827693-c2557780-ab70-11eb-90cd-b093004361db.gif" width="34%" />
    <img src="https://user-images.githubusercontent.com/696982/115995910-96a42180-a5b3-11eb-8200-1cfae06bf5bc.gif" width="34%" />
</p>

## Table of Contents

- Introduction (coming soon)
- [TODO](https://github.com/drgomesp/midgarts/blob/master/TODO.md#todo)
- [Dependencies](https://github.com/drgomesp/midgarts/blob/master/README.md#dependencies)
- [Building & Running](https://github.com/drgomesp/midgarts/blob/master/README.md#building-and-running)
- [Tools](https://github.com/drgomesp/midgarts/blob/master/README.md#tools)
    - [GRF Explorer](https://github.com/drgomesp/midgarts/blob/master/README.md#grf-explorer)
- [Examples](https://github.com/drgomesp/midgarts/blob/master/README.md#examples)

## Introduction

## TODO

Please have a look at the open milestones:

Milestone | Description |
--------- | ----------- |
[Character Graphics](https://github.com/drgomesp/midgarts/milestone/1) | Everything related to rendering character sprites, including character attachments, sprite animations and such.
[World Graphics](https://github.com/drgomesp/midgarts/milestone/2) | Everything related to world graphics, including 3D objects, terrain, water and lights.

## Dependencies

1. CentOS/Fedora-like Linux Distros:
   `SDL2{,_image,_mixer,_ttf,_gfx}-devel alsa-lib-devel libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel`

2. Arch Linux:
   `pacman -S sdl2{,_image,_mixer,_ttf,_gfx}`

3. MacOS:
  ```bash
  brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config
  ```

## Building and Running

1. Generate and env file by copying the distributed (.env.dist) file:
```bash
cp .env.dist .env 
```

2. Make sure to alter the `GRF_FILE_PATH` variable on the `.env` file:
```dotenv
GRF_FILE_PATH=/path/to/data.grf
```

3. Build the main binary by running:
```bash
go build -o midgarts ./cmd/sdlclient/main.go 
```

4. Run the binary:
```bash
./midgarts
```

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
