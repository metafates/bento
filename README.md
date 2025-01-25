# 🍱 Bento

> Work in progress. Mascot is too =)

<img width="1262" alt="image" src="https://github.com/user-attachments/assets/79a685a5-1168-496e-b386-462a04623494" />

<img align="right" width="250" src="https://github.com/user-attachments/assets/120b3882-0ed7-4b69-b5aa-2b75d83abe8a">

Bento is a Go framework for cooking up terminal user interfaces (TUIs). It provides a simple and flexible way to create text-based user interfaces in the terminal, which can be used for command-line applications, dashboards, and other interactive console programs.

It's a mix of Rust's [ratatui](https://ratatui.rs) library and Go's [bubbletea](https://github.com/charmbracelet/bubbletea)

Bento provides efficient and complex layout functionality from ratatui (widgets, buffers, cassowary constraint solving algorithm)
and functional design paradigms of [The Elm Architecture](https://guide.elm-lang.org/architecture/) as seen in bubbletea.

This library is mostly full Go rewrite of [ratatui-core](https://github.com/ratatui/ratatui/tree/main/ratatui-core) crate with some parts copied from bubbletea runtime.

Demo is WIP, but you can take a look at [examples](./examples) for now

> It's named bento after how similar some bento boxes look like compared to the typical TUIs (multiple blocks of different sizes side by side)
