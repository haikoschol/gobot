# gobot

## My first Go program

An IRC bot is my "hello world" for network applications. I have
[a few](https://github.com/haikoschol/dudebot)
[in Python](https://github.com/haikoschol/ircbotskel) and one in
[newLISP](https://github.com/haikoschol/repl_bot).

## Spec

The bot only supports connecting to a single server and channel.

The bot is started with three command line arguments:

* `server` the IRC server to connect to
* `channel` the channel to join
* `nick` the nickname to use

The bot understands one command: `msg`.

The `msg` command takes two arguments:

* the recipient of the message
* the content (body) of the message

There are public and private messages. A public message is sent by calling the
`msg` command in a channel. A private message is sent by calling the `msg`
command in a private chat with the bot.

There is a maximum number of messages per recipient. When that limit is reached,
no more messages are accepted and the bot responds with an error message.

There is a maximum number of recipients. When that limit is reached, no more
messages to new recipients are accepted and the bot responds with an error
message.