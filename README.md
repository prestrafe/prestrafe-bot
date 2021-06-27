## Building and running the bot locally

Start by checking out this repository and navigate into it within your terminal. Since the bot is designed to run inside
a Docker container, it requires some setup. First off, you will need to create a config file that contains the Twitch
channels that the bot should join. Create a file named `config.yml` in the source code directory and add the following
content:

```yaml
channels:
  - name: SomeChannel # This is the name of the Twitch channel that the bot should join. 
    gsiToken: xxx # This can be any random string, that needs to be present in your CSGO GSI config as well. 
    serverToken: xxx # Token string that will be registered in the server using sm_setprestrafetoken.
```

Next you will need to define some environment variables to configure the bots execution context. To do so, create a file
called `development.env` inside the source code directory and add the following content:

```properties
# This is required because of limitations of Docker for Windows, we cannot mount the config file.
BOT_CONFIGDIR=/src

# The API token the bot should used to authenticate against the Global API.
BOT_GLOBALAPITOKEN=xxx

# The address and port of the GSI backend service that should be used by the bot.
# It depends on how you run the GSI backend, but most likely example values are correct for local development.
BOT_GSIADDR=localhost
BOT_GSIPORT=8080

# The address and port of the SourceMod backend service that should be used by the bot.
BOT_SMADDR=localhost
BOT_SMPORT=8080

# The Twitch.tv username and API token that should be used to talk to the Twitch Chat API.
BOT_TWITCHUSERNAME=xxx
BOT_TWITCHAPITOKEN=xxx
```

Now you are ready to run the Docker container. The bot is shipped within a Docker container that builds the executable
and runs it directly afterwards. To build and run that container, perform the following commands:

```powershell
docker build -t prestrafe-bot:dev .
docker run --rm --network=host --name prestrafe-bot --env-file .\development.env -it prestrafe-bot:dev
```

Of course, you need to run the GSI backend service before the bot will be able to work.

## Supported commands

- `!bpb (bonus-number) (map-name)`: Displays the personal best time for the bonus stage.
- `!bwr (bonus-number) (map-name)`: Displays the world record time for the bonus stage.
- `!globalcheck`: Display global status of the server and player.
- `!prestrafe`: A list of supported commands of the Prestrafe bot.
- `!map (map-name)`: Displays information about the currently played map.
- `!mode`: Displays the currently played KZ timer mode.
- `!pb (map-name)`: Displays the personal best time for the main stage.
- `!rank (tp/pro/nub/all)`: Displays rank and map completion information.
- `!stats`: Displays a link to the GOKZ-Stats page.
- `!tier (map-name)`: Display the difficulty level for the map.
- `!wr (map-name)`: Displays the world record time for the main stage.
- `!server`: Displays server information (name, global status)
- `!run`: Displays current run information (Map name, course, checkpoints, teleports, time elapsed)

## Jumpstat Commands

- `!bh (bind|nobind)`: Displays the best bunny hop.
- `!dh (bind|nobind)`: Displays the best drop hop.
- `!laj (bind|nobind)`: Displays the best ladder jump.
- `!lj (bind|nobind)`: Displays the best long jump.
- `!mbh (bind|nobind)`: Displays the best multi bunny hop.
- `!wj (bind|nobind)`: Displays the best weird jump.

## Deployment Trigger

Number: 1
