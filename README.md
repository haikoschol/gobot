[![Build Status](https://travis-ci.org/haikoschol/gobot.svg?branch=master)](https://travis-ci.org/haikoschol/gobot)
[![Coverage Status](https://coveralls.io/repos/github/haikoschol/gobot/badge.svg?branch=master)](https://coveralls.io/github/haikoschol/gobot?branch=master)
[![Code Climate](https://codeclimate.com/github/haikoschol/gobot/badges/gpa.svg)](https://codeclimate.com/github/haikoschol/gobot)

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

Messages have the attributes
* sender
* timestamp
* private
* body

When a message to a new recipient is sent, the bot creates a mailbox for that
recipient. The mailbox is added to the repository of mailboxes.

The messages are stored until a user with the given recipient nick joins the
channel. Private messages are sent in a private chat, public ones in the
channel. They are sent in the format `<timestamp> <sender> <body>`. After a
message is sent, it is removed from the recipients mailbox. Empty mailboxes are
removed from the repository.

Messages are not sent immediately if the recipient is currently present in the
channel.

There is a maximum number of messages per mailbox. When that limit is reached,
no more messages to that recipient are accepted and the bot responds with an
error message.

There is a maximum number of mailboxes. When that limit is reached, no more
messages to new recipients are accepted and the bot responds with an error
message.
