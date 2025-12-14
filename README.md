<div align="center">

<pre>
,--------.,--.     ,--.        ,-----.,--.                      ,--. 
'--.  .--'|  ,---. `--' ,---. '  .--./|  ,---.  ,---. ,--.--. ,-|  | 
   |  |   |  .-.  |,--.(  .-' |  |    |  .-.  || .-. ||  .--'' .-. | 
   |  |   |  | |  ||  |.-'  `)'  '--'\|  | |  |' '-' '|  |   \ `-' | 
   `--'   `--' `--'`--'`----'  `-----'`--' `--' `---' `--'    `---'  
</pre>

</div>

## A Discord music bot hoping to act as a comparable replacement for jMusicBot, built with Go.
### This is meant to be self-hosted and can support 1 guild per instance.

## Summary

1. [Prerequisites](#prerequisites)
2. [Features](#features)
3. [Missing Features](#missing-features)
4. [Planned Additions](#planned-additions)
5. [Making a Discord Application](#making-a-discord-application)
   1. [General Information](#general-information)
   2. [Installation](#installation)
   3. [Bot](#bot)
6. [Getting Started with Lavalink](#getting-started-with-lavalink)
   1. [Installing Lavalink](#installing-lavalink)
   2. [Setting up Lavalink as a Service](#setting-up-lavalink-as-a-service)
   3. [Configuring Lavalink](#configuring-lavalink)
7. [Getting Started with ThisChord](#getting-started-with-thischord)
   1. [Installing ThisChord](#installing-thischord)
   2. [Setting up ThisChord as a Service](#setting-up-thischord-as-a-service)
   3. [Configuring ThisChord](#configuring-thischord)
8. [Starting the Music Bot](#starting-the-music-bot)

### Prerequisites

* An application created in the [Discord Developer Portal](https://discord.com/developers/applications)
* Go 1.25.1+
* Java 17 or higher
* Lavalink
* [yt-dlp](https://github.com/yt-dlp/yt-dlp/wiki/Installation#third-party-package-managers)
* [ffmpeg](https://www.ffmpeg.org/download.html)

### Features

* Search for music on YouTube
* Play music from YouTube by search query or URI
* Add additional songs into a queue
* Skip songs
* Shuffle the queue
* Clear the queue
* See what's playing, and the position of the seeker
* Stop & pause playback

### Missing Features

* Ability to repeat:
  * tracks
  * queues
* Loading playlists from YouTube
* Additional music sources:
  * SoundCloud
  * Bandcamp
  * Vimeo
  * Premium services like:
    * Spotify
    * Apple Music
    * Deezer

### Planned Additions

* YAML or JSON config instead of using .env file
* Ability to run w/ Docker
* CI/CD integration
* Return custom embeds instead of plain old text responses:
  * Now Playing:
    * media control buttons
    * seeker progression
    * added by
    * art (cover/album art)
  * Searching
  * Skipping
  * etc.

## Making a Discord Application

First, we want to navigate to the [Discord Developer Portal](https://discord.com/developers/applications) and click **New Application** in the top right, then proceed with editing the following sections.

### General Information

We don't really need to change anything here unless you decide to change the name of the application.

### Installation

Check **Guild Install** only, and ensure the scope has `applications.commands` and `bot` in it, then add the following permissions:

* Add Reactions
* Connect
* Embed Links
* Priority Speaker
* Read Message History
* Request To Speak
* Send Messages
* Send Messages in Threads
* Speak
* Use Slash Commands
* Use Voice Activity
* View Channels

> [!TIP]
> You'll use the URL found under `Install Link` to invite the bot to your server.

### Bot

Here, give your bot an optional profile picture, and a username. Ensure the following are enabled:

* Public Bot
* Presence Intent
* Server Members Intent

> [!TIP]
> This is where you'll also get your bot token, under the `Token` section of this page.

## Getting Started with Lavalink

### Installing Lavalink

Please follow the official Lavalink documentation [here](https://lavalink.dev/getting-started/binary.html)

### Setting up Lavalink as a Service

Please follow the official Lavalink documentation on creating a systemd service [here](https://lavalink.dev/getting-started/systemd.html)

### Configuring Lavalink

Lavalink needs a bit of additional configuration before we can use it with YouTube. You'll want to create a file called `application.yml` in the directory with `Lavalink.jar` - I recommend copying the [sample](https://github.com/eleinah/thischord/blob/main/sample-application.yml) and changing the password to your liking.
 
## Getting Started with ThisChord 

### Installing ThisChord

> [!IMPORTANT]
> Please ensure you have Go 1.25.1+ on your machine; installation instructions can be found [here](https://go.dev/doc/install)

To install ThisChord, run the following in your terminal, which will output the binary file to `$GOBIN`

```sh
$ go install github.com/eleinah/thischord/cmd/thischord@latest
```

> [!TIP]
> Either export `GOBIN` in your shell profile, or run `go env -w GOBIN=/path/to/some/bin`

### Setting up ThisChord as a Service

Create a file in `/etc/systemd/system` named `thischord.service` with the following:

```
[Unit]
Description=ThisChord Music Bot
After=lavalink.service syslog.target network.target
Requires=lavalink.service

[Service]
User=BOT_USER # replace this with the user you run the bot as, i.e. bot
Group=BOT_GROUP # replace this with the user's group you run the bot as, i.e., bot
WorkingDirectory=/path/to/bin/with/bot # i.e. /home/bot
ExecStart=/path/to/bin/with/bot/thischord # i.e /home/bot/thischord
Restart=on-failure
RestartSec=5s


[Install]
WantedBy=multi-user.target
```

### Configuring ThisChord

Navigate to `$GOBIN`, or wherever you placed the `thischord` binary, and create a `.env` file. I recommend copying the [sample](https://github.com/eleinah/thischord/blob/main/.env.sample) and filling it in with details specific to your instance and Discord guild.

## Starting the Music Bot

> [!IMPORTANT]
> Please ensure you have [yt-dlp](https://github.com/yt-dlp/yt-dlp/wiki/Installation#third-party-package-managers) and [ffmpeg](https://www.ffmpeg.org/download.html) installed before starting the bot.

You'll want to reload the systemctl daemon:

```sh
$ sudo systemctl daemon-reload
```

and then enable & start both the Lavalink and ThisChord services:

```sh
$ sudo systemctl enable --now lavalink.service
$ sudo systemctl enable --now thischord.service
```

Lastly, let's make sure things are running smoothly. Run the following commands to ensure everything started properly:

```sh
$ systemctl status lavalink.service
$ systemctl status thischord.service
```

That should be all! You now have a music bot for your server! Make sure to periodically check this repo for updates!

