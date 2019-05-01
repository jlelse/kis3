# KISSS (Keep It Simple Stupid Stats) a.k.a kis3

> I started to develop KISSS because I need a simple and privacy respecting method to collect visitor statistics of my website. I don't need any fancy dashboard, but I want to get exactly the stats I need. I need something fast, that is able to run even on really low end hardware like a Raspberry Pi.

## How to install

KISSS is really easy to install via Docker.

    docker run -d --name kis3 -p 8080:8080 -v kis3:/app/data -v ${pwd}/config.json:/app/config.json kis3/kis3

Depending on your setup, replace `-p 8080:8080` with your custom port configuration. KISSS listens to port 8080 by default, but you can also change this via the configuration.

To persist the data KISSS collects, you should mount a volume or a folder to `/app/data`. When mounting an folder, give writing permissions to UID 100, because it is a non-root image to make it more secure.

You should also mount a configuration file to `/app/config.json`.

It's also possible to use KISSS without Docker, but for that you need to compile it yourself. In the future there will be executables without dependencies available.

## Configuration

You can configure some settings using a `config.json` file in the working directory or provide a custom path to it using the `-c` CLI flag:

`port` (`8080`): Set the port to which KISSS should listen

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

If you specify an environment variable (`PORT`, `DNT`, `DB_PATH`, `STATS_USERNAME`, `STATS_PASSWORD`), that will override the settings from the configuration file.

## Add to website

You can add the KISSS tracker to any website by putting `<script src="https://yourkis3domain.tld/kis3.js"></script>` just before `</body>` in the HTML. Just replace `yourkis3domain.tld` with the correct address.

## Requesting stats

You can request statistics via the `/stats` endpoint and specifying filters via query parameters (`?view=hours&format=chart` etc.). By combining this filters, you can exactly request the stats you want to get.

The following filters are available:

`view`: specify what data (and it's view counts) gets displayed, you have the option between `pages` (tracked URLS), `referrers` (tracked refererrers - only hostnames e.g. google.com), `useragents` (tracked useragents with version - browsers or crawl bots with version), `useragentnames` (tracked useragents without version), `hours` / `days` / `weeks` / `months` (tracks grouped by hours / days / weeks / months), `allhours` / `alldays` (tracks grouped by hours / days including hours or days with zero visits, spanning from first to last track in selection)

`from`: start time of the selection in the format `YYYY-MM-DD HH:MM`, e.g. `2019-01` or `2019-01-01 01:00`

`to`: end time of the selection

`url`: filter URLs containing the string provided, so `word` filters out all URLs that don't contain `word`

`ref`: filter referrers containing the string provided, so `word` filters out all refferers that don't contain `word`

`ua`: filter user agents containing the string provided, so `Firefox` filters out all user agents that don't contain `Firefox`

`orderrow`: row to use for ordering, `first` for the data groups, `second` for the view counts

`order`: select whether to use ascending order (`ASC`) or descending order (`DESC`)

`format`: the format to represent the data, default is `plain` for a simple plain text list, `json` for a JSON response or `chart` for a chart generated with ChartJS in the browser

## Daily email reports

KISSS has a feature that can send you daily email reports. It basically requests the statistics and sends the response via email. You can configure it by adding report configurations to the configuration file:

```json
{
  // Other configurations...
  "reports": [
    {
      "name": "Daily stats from KISSS",
      "time": "15:00",
      "query": "view=pages&orderrow=second&order=desc",
      "from": "myemailaddress@mydomain.tld",
      "to": "myemailaddress@mydomain.tld",
      "smtpHost": "mail.mydomain.tld:587",
      "smtpUser": "myemailaddress@mydomain.tld",
      "smtpPassword": "mysecretpassword"
    },
    {
      // Additional reports...
    }
  ]
}
```

## License

KISSS is licensed under the MIT license, so you can do basically everything with it, but nevertheless, please contribute your improvements to make KISSS better for everyone. See the LICENSE.txt file.
