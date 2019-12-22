# gossh

ssh to websocket, webssh backend

通过websocket 转ssh 链接其他主机

Ssh via websocket to other hosts

做为webssh后端,需与xteam.js 配合使用

As a webssh backend, it needs to work with xteam.js

### Install xteam.js
```bash
npm install --save xterm.js 
npm install --save xterm-addon-attach
```

### Use

```vue
<template>
<div ref="terminal" id="terminal"></div>
</template>
import { AttachAddon } from 'xterm-addon-attach'
import { Terminal } from 'xterm'
import 'xterm/css/xterm.css'

export default {
  name: 'CompTerm',
  data () {
    return {
      isFullScreen: false,
      searchKey: '',
      v: false,
      ws: null,
      term: null,
        hostname: '0.0.0.0',
        Port: 22,
        username: 'root',
        password: '',
        pri_key: '',
        wsUrl: 'ws://127.0.0.1:9018/ws?' + 'cols=' + this.term.cols + '&rows=' + this.term.rows + '&hostname=' + this.hostname + '&port=' + this.Port + '&user=' + this.username + '&pd=' + this.password + '&pk=' + this.pri_key
    }
  }

this.term = new Terminal({
        rows: 35,
        fontSize: 18,
        cursorBlink: true,
        cursorStyle: 'bar',
        bellStyle: 'sound',
        theme: defaultTheme
      })
 this.term.open(this.$refs.terminal)
 this.ws = new WebSocket(this.wsUrl)
const attachAddon = new AttachAddon(this.ws)
this.term.loadAddon(attachAddon)
```

## parameter

```golang
	Hostname string `form:"hostname" json:"hostname" binding:"required"`
	Port     string `form:"port,default=22" json:"port" `
	User     string `form:"user,default=root" json:"user"`
	Password string `form:"pd" json:"password"` //base64 encoded password
	Cols     int    `form:"cols,default=120" json:"cols"`
	Rows     int    `form:"rows,default=32" json:"rows"`
	PriKey   string `form:"pk" json:"pri_key"` //base64 private key
```
## Build 
```bash
go build 
```
