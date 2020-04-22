# KISSS (Keep It Simple Stupid Stats) a.k.a kis3

> I started to develop KISSS because I need a simple and privacy respecting method to collect visitor statistics of my website. I don't need any fancy dashboard, but I want to get exactly the stats I need. I need something fast, that is able to run even on really low end hardware like a Raspberry Pi.

## How to install

KISSS is really easy to install via Docker.

```bash
docker run -d --name kis3 -p 8080:8080 -v kis3:/app/data -v ${pwd}/config.json:/app/config.json kis3/kis3
```

Depending on your setup, replace `-p 8080:8080` with your custom port configuration. KISSS listens to port 8080 by default, but you can also change this via the configuration.

To persist the data KISSS collects, you should mount a volume or a folder to `/app/data`. When mounting an folder, give writing permissions to UID 100, because it is a non-root image to make it more secure.

You should also mount a configuration file to `/app/config.json`.

### Build from source

It's also possible to use KISSS without Docker, but for that you need to compile it yourself. All you need to do so is installing go (follow the [instruction](https://golang.org/doc/install) or use [distro.tools](https://distro.tools) to install the latest version on Linux - you need at least version 1.14) and execute the following command:

```bash
go get -u kis3.dev/kis3
```

After that there should be an executable with the name `kis3` in `$HOME/go/bin`.

## Configuration

You can configure some settings using a `config.json` file in the working directory or provide a custom path to it using the `-c` CLI flag:

`port` (`8080`): Set the port to which KISSS should listen

`baseUrl` (optional, required for Telegram): Set the base URL on which KISSS runs

`dnt` (`true`): Set whether or not KISSS should respect Do-Not-Track headers some browsers send

`dbPath` (`data/kis3.db`): Set the path for the SQLite database (relative to the working directory - in the Docker container it's `/app`).

You can make the statistics private and only accessible with authentication by setting both `statsUsername` and `statsPassword` to a username and password. If only one or none is set, the statistics are accessible without authorization and public to anyone.

The configuration file can look like this:

```json
{
  "port": 8080,
  "dnt": true,
  "dbPath": "data/kis3.db",
  "statsUsername": "myusername",
  "statsPassword": "mysecretpassword"
}
```

If you specify an environment variable (`PORT`, `BASE_URL`, `DNT`, `DB_PATH`, `STATS_USERNAME`, `STATS_PASSWORD`), that will override the settings from the configuration file.

### Email

To enable email integration for sending reports, you need to add some configuration values for that:

`smtpFrom`: Sender address for the emails

`smtpHost`: Address of the mail server (including port)

`smtpUser`: Username for SMTP login

`smtpPassword`: Password for SMTP login

### Telegram

The Telegram integration allows sending reports via Telegram and also requesting stats via Telegram. For that the following configuration value must be set:

`tgBotToken`: Token for the Telegram bot, which you can request via the [Bot Father](https://t.me/BotFather)

`tgHookSecret` (optional): Secret, so nobody (except KISSS and Telegram) knows the URL to which Telegram should send updates about new messages

## Add to website

You can add the KISSS tracker to any website by putting `<script src="https://yourkis3domain.tld/kis3.js"></script>` just before `</body>` in the HTML. Just replace `yourkis3domain.tld` with the correct address.

## Requesting stats

You can request statistics via the `/stats` endpoint and specifying filters via query parameters (`?view=hours&format=chart` etc.). By combining this filters, you can exactly request the stats you want to get.

The following filters are available:

`view`: specify what data (and it's view counts) gets displayed, you have the option between `pages` (tracked URLS), `referrers` (tracked refererrers - only hostnames e.g. google.com), `useragents` (tracked useragents with version - browsers or crawl bots with version), `useragentnames` (tracked useragents without version), `os` (tracked operating systems), `hours` / `days` / `weeks` / `months` (tracks grouped by hours / days / weeks / months), `allhours` / `alldays` (tracks grouped by hours / days including hours or days with zero visits, spanning from first to last track in selection), `count` (count all tracked views where filters apply)

`from`: start time of the selection in the format `YYYY-MM-DD HH:MM`, e.g. `2019-01` or `2019-01-01 01:00`

`to`: end time of the selection

`fromrel` / `torel`: relative time from now for `to` or `from` (e.g `-2h45m`, valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"

`url`: filter URLs containing the string provided, so `word` filters out all URLs that don't contain `word`

`ref`: filter referrers containing the string provided, so `word` filters out all refferers that don't contain `word`

`ua`: filter user agents containing the string provided, so `Firefox` filters out all user agents that don't contain `Firefox`

`os`: filter operating systems containing the string provided, so `Windows` filters out all operating systems that don't contain `Windows`

`bots`: filter out bots (`0`) or show only bots (`1`)

`ordercol`: column to use for ordering, `first` for the data groups, `second` for the view counts

`order`: select whether to use ascending order (`ASC`) or descending order (`DESC`)

`limit`: limit the number of rows returned

`format`: the format to represent the data, default is `plain` for a simple plain text list, `json` for a JSON response or `chart` for a chart generated with ChartJS in the browser

### Via Telegram

You can also request stats via Telegram (in case you enable the Telegram integration). To do so, simply send a message with the command `stats` like `/stats view=pages...`.

If you have authentication enabled, you need to add `username=yourusername&password=yourpassword` to the query.

## Daily reports

KISSS has a feature that can send you daily reports. It basically requests the statistics and sends the response via your preferred communication channel (mail or Telegram). You can configure it by adding report configurations to the configuration file:

```json
{
  // Other configurations...
  "reports": [
    {
      // Email configuration
      "name": "Daily stats from KISSS",
      "time": "15:00",
      "query": "view=pages&bots=0&ordercol=second&order=desc",
      "from": "myemailaddress@mydomain.tld",
      "to": "myemailaddress@mydomain.tld"
    },
    {
      // Telegram configuration
      "name": "Daily stats from KISSS",
      "type": "telegram", // Add this for Telegram
      "time": "15:00",
      "query": "view=pages&bots=0&ordercol=second&order=desc",
      "tgUserId": 123456
    },
    {
      // Additional reports...
    }
  ]
}
```

You can find out your Telegram user id using [@userinfobot](https://t.me/userinfobot).

## License

KISSS is licensed under the MIT license, so you can do basically everything with it, but nevertheless, please contribute your improvements to make KISSS better for everyone. See the LICENSE.txt file.
