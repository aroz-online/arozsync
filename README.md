# arozsync
Synchronize your arozos folders across different PCs and laptops (WIP)

## Installation

### Requirements

- Go 1.17+
- npm (Node 15+)
- Wails

This project require wails for generating the native UI for different platform. You can get it [here](https://wails.io/docs/gettingstarted/installation).



### Build From Source

For more information, read wail's build instruction

#### Development build

```
wails dev
```

#### Production build

```
wails build
```





## Compatibility Issues

If you encountered any compatibility issues related to the webview / windows interface like "app not starting", "it crashed during start" or "it flash once and dissappeared", make sure you have [Microsoft Edge Webview 2](https://developer.microsoft.com/en-US/microsoft-edge/webview2/) installed. If the issue still not resolved after updating your webview2 element, raise your issue on the [go-webview2](https://github.com/jchv/go-webview2) library side instead.

## License

MIT License