# redcaster

A ray caster built from scratch.

## Inspiration

Being lucky enough to have been exposed to computers growing up in the 80s and 90s, I had a front row seat for the
releases Wolfenstein 3D, DOOM, and later Quake.

I was exposed to John Carmack's graphics wizardry: build immersive worlds with limited hardware. Every time id Software
released something, you knew it would blow away everything that came before.

Over the years, after reading "Masters of Doom" a few times, reading John Romero's "Doom Guy", and Fabien Sanglard's
game engine black books (highly recommend them, btw), I got nostalgic and was somewhat was curious about how it would
have felt to be a game dev back in those days. That's why this project got started in the first place.

## The Challenge

Although the hardware limitations are nowhere near what they used to be back then, I wanted to write something
completely from scratch that only leverages only the CPU. Pure software rendering. No GPU assistance, no 3D library, no
game engine, nothing. Oh, and in Golang on top of it.

Why Golang? it's the language I've been the most excited about for the past few years honestly. Although not optimal or 
low level as C/C++ or Rust, I thought it would be an interesting challenge to optimize for and avoid garbage collection.
I was genuinely intrigued to see if a GCed language would prove to be a worthy contender for building such a project.

Since it's from scratch, I intentionally minimized dependencies to almost nothing. Don't believe me? have a look at
the `go.mod` file.

The only dependency is gofw for creating the window, and poll input. Everything is rendered to a byte array serving as a
frame buffer by setting RGBA values for each individual pixel, at each frame.

The `redpix` library has a simple interface to provide the byte array framebuffer and it will make sure it gets
displayed on screen.

## What's next

This is a WIP project. It's been sitting for about a year, and I'm picking it back up in 2026 (new inspiration from
DOOM).
