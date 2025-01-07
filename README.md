# Bento 🍱

> Work in progress

Bento is a Go library for cooking up delicious TUIs (terminal user interfaces).

It's a mix of Rust's [ratatui](https://ratatui.rs) library and Go's [bubbletea](https://github.com/charmbracelet/bubbletea)

Bento provides efficient and complex layout functionality from ratatui (widgets, buffers, cassowary constraint solving algorithm)
and functional design paradigms of [The Elm Architecture](https://guide.elm-lang.org/architecture/) as seen in bubbletea.

This library is mostly full Go rewrite of [ratatui-core](https://github.com/ratatui/ratatui/tree/main/ratatui-core) crate with some parts copied from bubbletea runtime.

No demo yet, but you can take a look at [examples](./examples) for now

> It's named bento after how similar some bento boxes look like compared to the typical TUIs (multiple blocks of different sizes side by side)
